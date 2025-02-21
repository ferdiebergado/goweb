import { FormHandler } from './form_handler';

const regFrm = new FormHandler({
  validationRules: {
    email: (value) => {
      if (!value) return 'Email is required';
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
      return emailRegex.test(value) ? undefined : 'Invalid email format';
    },
    password: (value) => (value ? undefined : 'Password is required'),
    password_confirm: (value) =>
      value ? undefined : 'password_confirm is required',
  },
  onSuccess: (data) => {
    console.log('Success:', data);
    alert('Form submitted successfully!');
  },
  onError: (error) => {
    console.error('Error:', error);
    if (typeof error === 'object') {
      for (const key in error) {
        alert(`${key}: ${error[key]}`); // Display specific error messages
      }
    } else {
      alert(error); // Or display a general error message
    }
  },
  onFinally: () => {
    console.log('Form submission process finished.');
  },
});

regFrm.handleSubmit();
