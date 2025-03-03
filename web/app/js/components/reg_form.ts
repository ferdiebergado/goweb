import { FormErrors, FormValues } from '../@types/form';
import { isValidEmail } from '../utils';
import form from './form';

interface RegData extends FormValues {
  email: string;
  password: string;
  password_confirm: string;
}

interface RegErrors extends FormErrors {
  email?: string;
  password?: string;
  passwordConfirm?: string;
}

export default function () {
  return form<RegData, RegErrors>({
    data: {
      email: '',
      password: '',
      password_confirm: '',
    },
    submitUrl: '/api/auth/register',
    errors: {},
    validateFn() {
      const { email, password, password_confirm } = this.data;
      const errors: RegErrors = {};

      if (!email) {
        errors.email = 'Email is required.';
      } else if (!isValidEmail(email)) {
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
    onError(error) {
      console.error(error);
    },
  });
}
