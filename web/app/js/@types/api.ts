export interface APIResponse<T> {
  message: string;
  errors: Record<string, string | undefined>;
  data: T;
}
