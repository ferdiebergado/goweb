import { APIResponse } from './api';

interface FormOptions {
  data: Record<string, string | number | boolean>;
  method?: 'POST' | 'PUT' | 'PATCH' | 'DELETE';
  submitUrl: string;
  errors: Record<string, string | undefined>;
  validateFn(): Record<string, string | undefined>;
  onSuccess<K>(data: APIResponse<K>): undefined;
  onError(error: unknown): undefined;
}

export { FormOptions };
