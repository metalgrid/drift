import "@girs/gjs";

import $t from "@girs/st-15";
import {
  Extension,
  gettext as _,
} from "@girs/gnome-shell/extensions/extension";
import * as PanelMenu from "@girs/gnome-shell/ui/panelMenu";
import * as Main from "@girs/gnome-shell/ui/main";

export default class Drift extends Extension {
  private _indicator: PanelMenu.Button | null = null;

  override enable(): void {
    this._indicator = new PanelMenu.Button(0, _(this.metadata.name), false);

    const icon = new $t.Icon({
      iconName: "send-to-sybolic",
      styleClass: "system-status-icon",
    });
    this._indicator.add_child(icon);

    Main.panel.addToStatusArea(this.uuid, this._indicator);
  }

  override disable(): void {
    this._indicator?.destroy();
    this._indicator = null;
  }
}
