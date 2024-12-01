import { ButtonProps, DriftService } from "./types";
import "@girs/gjs";
import "@girs/gjs/dom";

import $t from "@girs/st-15";
import Gio from "@girs/gio-2.0";
import GLib from "@girs/glib-2.0";
import GObject from "@girs/gobject-2.0";
import {
  Extension,
  gettext as _,
} from "@girs/gnome-shell/extensions/extension";
import "@girs/gnome-shell/extensions/global";
import * as Main from "@girs/gnome-shell/ui/main";
import {
  Notification,
  NotificationApplicationPolicy,
  NotificationDestroyedReason,
  Source,
} from "@girs/gnome-shell/ui/messageTray";
import * as PanelMenu from "@girs/gnome-shell/ui/panelMenu";
import { PopupMenu, PopupMenuItem } from "@girs/gnome-shell/ui/popupMenu";
const serviceName = "com.github.metalgrid.Drift";
const objPath = "/com/github/metalgrid/Drift";
const ifaceName = "com.github.metalgrid.Drift";
const sigNotify = "Notify";
const sigQuestion = "Question";
class PersistentNotifierType extends Source {
  public override destroy(_reason: NotificationDestroyedReason): void {}

  public sendNotification(title: string, body: string, icon?: string): void {
    const notification = new Notification({
      source: this,
      title: title,
      body: body,
    });

    if (icon) {
      notification.gicon = new Gio.ThemedIcon({ name: icon });
    }

    this.addNotification(notification);
  }

  public sendActionableNotification(
    title: string,
    body: string,
    actions: { [label: string]: () => void }
  ): void {
    const notification = new Notification({
      source: this,
      title: title,
      body: body,
    });

    for (const [label, callback] of Object.entries(actions)) {
      notification.addAction(label, () => {
        callback();
      });
    }

    this.addNotification(notification);
  }

  public _destroy(_reason: NotificationDestroyedReason): void {}

  public real_destroy(reason: NotificationDestroyedReason): void {
    return super.destroy(reason);
  }
}

const PersistentNotifier = GObject.registerClass(PersistentNotifierType);

function createProxy(
  bus: Gio.DBusConnection,
  service: string,
  path: string,
  schema: string
): DriftService {
  const DriftProxyFactory = Gio.DBusProxy.makeProxyWrapper(schema);
  return new (DriftProxyFactory as any)(bus, service, path);
}

const Drift = GObject.registerClass(
  class Drift extends PanelMenu.Button {
    private driftService: DriftService;
    private bus = Gio.DBus.session;
    public extDir: Gio.File;
    private notifier: PersistentNotifierType;
    private subscriptions: number[] = [];
    private icon: $t.Icon;

    public constructor(path: Gio.File, args: ButtonProps) {
      //@ts-expect-error // Smells like a TS bug. ButtonProps *is* a union of tuples.
      super(...args);
      this.extDir = path;

      const theme = $t.ThemeContext.get_for_stage(
        global.get_stage()
      ).get_theme();
      theme.load_stylesheet(
        Gio.file_new_for_path(`${this.extDir.get_path()}/stylesheet.css`)
      );

      const [ok, buf] = GLib.file_get_contents(
        `${this.extDir.get_path()}/schema/${ifaceName}.xml`
      );
      if (!ok) {
        logError(
          new Error(
            `Failed to load schema for ${ifaceName} from ${this.extDir}`
          )
        );
        throw new Error(`Failed to load schema for ${ifaceName}`);
      }
      this.driftService = createProxy(
        Gio.DBus.session,
        ifaceName,
        objPath,
        new TextDecoder().decode(buf)
      );

      this.notifier = new PersistentNotifier({
        title: _("Drift"),
        icon: new Gio.ThemedIcon({ name: "mail-unread-symbolic" }),
        policy: new NotificationApplicationPolicy(serviceName),
      });
      Main.messageTray.add(this.notifier);

      this.icon = new $t.Icon({
        iconName: "drift-symbolic",
        styleClass: "system-status-icon",
      });
      this.add_child(this.icon);

      this.bus.watch_name(
        serviceName,
        Gio.BusNameWatcherFlags.NONE,
        this.driftOnline.bind(this),
        this.driftOffline.bind(this)
      );
    }

    driftOnline() {
      let subId = this.bus.signal_subscribe(
        serviceName,
        ifaceName,
        sigNotify,
        objPath,
        null,
        Gio.DBusSignalFlags.NONE,
        (
          _connection,
          _senderName,
          _objectPath,
          _interfaceName,
          _signalName,
          parameters
        ) => {
          log("Handling notify signal");
          const [message] = parameters.recursiveUnpack();
          this.notifier.sendNotification("Drift", message);
        }
      );
      this.subscriptions.push(subId);

      subId = this.bus.signal_subscribe(
        serviceName,
        ifaceName,
        sigQuestion,
        objPath,
        null,
        Gio.DBusSignalFlags.NONE,
        (
          _connection,
          _senderName,
          _objectPath,
          _interfaceName,
          _signalName,
          parameters
        ) => {
          const [id, message] = parameters.recursiveUnpack();

          this.notifier.sendActionableNotification("Drift", message, {
            Accept: () => {
              this.driftService.RespondSync(id, "ACCEPT");
            },
            Decline: () => {
              this.driftService.RespondSync(id, "DECLINE");
            },
          });
        }
      );
      this.subscriptions.push(subId);

      this.connect("button-press-event", () => {
        // TS doesn't know that this.menu is always a PopupMenu
        this.menu = this.menu as PopupMenu;
        const [peers] = this.driftService.ListPeersSync();
        log(`${peers.length} peers found`);
        this.menu.removeAll();
        peers.forEach((peer) => {
          const menuItem = new PopupMenuItem(peer);
          menuItem.connect("activate", () => {
            this.selectFile().then((uris) => {
              const file = Gio.File.new_for_uri(uris[0]!);
              this.driftService.RequestSync(peer, file.get_path()!);
            });
          });
          (this.menu as PopupMenu).addMenuItem(menuItem);
        });

        this.menu.open(true);
      });

      this.icon.styleClass = "system-status-icon icon-online";
    }

    // This uses the freedesktop portal to show a file selection dialog.
    // Probably not the cleanest way, but I can't figure out a simpler convenient way.
    selectFile() {
      return new Promise<string[]>((resolve, reject) => {
        this.bus
          .call(
            "org.freedesktop.portal.Desktop",
            "/org/freedesktop/portal/desktop",
            "org.freedesktop.portal.FileChooser",
            "OpenFile",
            GLib.Variant.new_tuple([
              GLib.Variant.new_string(""),
              GLib.Variant.new_string(_("Send file")),
              new GLib.Variant("a{sv}", {
                handle_token: GLib.Variant.new_string(`drift_${Date.now()}`),
                accept_label: GLib.Variant.new_string(_("Send")),
              }),
            ]),
            null,
            null,
            -1,
            null
          )
          .then((result) => {
            const [handle] = result.deepUnpack<[string]>();

            this.bus.signal_subscribe(
              "org.freedesktop.portal.Desktop",
              "org.freedesktop.portal.Request",
              "Response",
              handle,
              null,
              Gio.DBusSignalFlags.NONE,
              (_sender, _iface, _path, _name, _signal, params) => {
                const [_, { uris }] = params.recursiveUnpack() as any;
                resolve(uris);
              }
            );
          })
          .catch((e) => reject(e));
      });
    }

    driftOffline() {
      log("DBus service offline");
      this.subscriptions.forEach((subId) => {
        this.bus.signal_unsubscribe(subId);
      });
      this.icon.styleClass = "system-status-icon icon-offline";
    }

    override destroy(): void {
      if (this.notifier) {
        this.notifier._destroy(NotificationDestroyedReason.SOURCE_CLOSED);
      }
    }
  }
);

export default class DriftExtension extends Extension {
  private extension: any = null;

  override enable(): void {
    this.extension = new Drift(this.metadata.dir, [
      0,
      _(this.metadata.name),
      false,
    ]);
    Main.panel.addToStatusArea(this.uuid, this.extension);
  }

  override disable(): void {
    this.extension?.destroy();
    this.extension = null;
  }
}
