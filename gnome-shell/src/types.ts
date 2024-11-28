// Type definitions for the Drift service
import { Button } from "@girs/gnome-shell/ui/panelMenu";

export type ButtonProps =
  | [Partial<Button.ConstructorProps>]
  | [menuAlignment: number, nameText: string, dontCreateMenu?: boolean];

export interface DriftService {
  ListPeersSync(): [string[]];
  ListPeersAsync(): Promise<[string[]]>;
  RequestSync(peer: string, file: string): void;
  RequestAsync(peer: string, file: string): Promise<void>;
  RespondSync(id: string, answer: "ACCEPT" | "DECLINE"): void;
  RespondAsync(id: string, answer: "ACCEPT" | "DECLINE"): Promise<void>;
}
