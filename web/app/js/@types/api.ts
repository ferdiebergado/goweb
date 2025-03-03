import { FormErrors } from './form';

export interface APIResponse<T> {
  message: string;
  errors: FormErrors;
  data: T;
}
