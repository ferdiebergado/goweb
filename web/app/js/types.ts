export interface APIResponse<T> {
  message: string;
  errors: { [key: string]: string };
  data: T;
}
