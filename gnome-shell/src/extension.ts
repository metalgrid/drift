import { DriftService } from "./types";
import "@girs/gjs";

import $t from "@girs/st-15";
import Gio from "@girs/gio-2.0";
import GObject from "@girs/gobject-2.0";

import {
  Extension,
  gettext as _,
} from "@girs/gnome-shell/extensions/extension";
import * as PanelMenu from "@girs/gnome-shell/ui/panelMenu";
import * as Main from "@girs/gnome-shell/ui/main";

const serviceName = "com.github.metalgrid.Drift";
const objPath = "/com/github/metalgrid/Drift";
const ifaceName = "com.github.metalgrid.Drift";

const Drift = GObject.registerClass(
  class Drift extends PanelMenu.Button {
    private _driftProxy: DriftService | null = null;
    // nasty hacks
    public _path: string = "";

    override _init(
      params?: Partial<PanelMenu.ButtonBox.ConstructorProps> & { path: string }
    ): void;

    override _init(
      menuAlignment: number,
      nameText: string,
      dontCreateMenu?: boolean
    ): void;

    override _init(
      _p1?: Partial<PanelMenu.ButtonBox.ConstructorProps> | number,
      _p2?: string,
      _p3?: boolean
    ): void {
      super._init(...arguments);
      this.init();
    }

    init() {
      const icon = new $t.Icon({
        iconName: "send-to-symbolic",
        styleClass: "system-status-icon",
      });

      this.add_child(icon);
      this.setupDBus();
    }

    setupDBus() {
      const bus = Gio.DBus.session;

      bus.watch_name(
        serviceName,
        Gio.BusNameWatcherFlags.NONE,
        () => this._serviceOnline(),
        () => this._serviceOffline()
      );
    }

    _serviceOnline() {
      log("DBus service online");

      // should be read from the bundled schema...
      const schema = `
      <node>
    <interface name="com.github.metalgrid.Drift">
        <method name="ListPeers">
            <arg type="as" direction="out"></arg>
        </method>
        <method name="Request">
            <arg type="s" direction="in"></arg>
            <arg type="s" direction="in"></arg>
        </method>
    </interface>
</node>
      `;

      const DriftProxy = Gio.DBusProxy.makeProxyWrapper(schema);
      this._driftProxy = new (DriftProxy as any)(
        Gio.DBus.session,
        serviceName,
        objPath
      );
      const peers = this._driftProxy?.ListPeersSync();
      log("Drifters: ", peers);
    }

    _serviceOffline() {
      log("DBus service offline");
    }
  }
);

export default class DriftExtension extends Extension {
  private extension: any = null;

  override enable(): void {
    this.extension = new Drift(0, _(this.metadata.name), false);
    this.extension._path = this.metadata.path;
    Main.panel.addToStatusArea(this.uuid, this.extension);
  }

  override disable(): void {
    this.extension?.destroy();
    this.extension = null;
  }
}
