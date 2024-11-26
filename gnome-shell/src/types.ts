export interface DriftService {
  ListPeersSync(): [string];
  ListPeersAsync(): Promise<[string]>;
  RequestSync(peer: string, file: string): void;
  RequestAsync(peer: string, file: string): Promise<void>;
}
