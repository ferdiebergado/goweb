import { FormErrors } from '../@types/form';
import { isValidEmail } from '../utils';
import form from './form';

type FormValues = {
  email: string;
  password: string;
  passwordConfirm: string;
};

export default function () {
  return form({
    data: {} as FormValues,
    submitUrl: '/api/auth/register',
    errors: {} as FormErrors<FormValues>,
    validateFn() {
      const { email, password, password_confirm } = this.data;
      const errors: FormErrors<FormValues> = {};

      if (!email) {
        errors.email = 'Email is required.';
      } else if (!isValidEmail(email as string)) {
        errors.email = 'Invalid email format.';
      }

      if (!password) {
        errors.password = 'Password is required.';
      }

      if (!password_confirm) {
        errors.passwordConfirm = 'Password confirmation is required.';
      } else if (password && password_confirm !== password) {
        errors.passwordConfirm = 'Passwords should match.';
      }

      return errors;
    },
    onSuccess(res) {
      const { message, data } = res;
      console.log(message, data);
    },
    onError() {
      return;
    },
  });
}
