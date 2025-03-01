interface LoginData {
  email: string;
  password: string;
  loading: boolean;
  message: string;
  errors: { email?: string; password?: string };
  validate(): boolean;
  register(): Promise<void>;
}

export default (): LoginData => ({
  email: '',
  password: '',
  loading: false,
  message: '',
  errors: {},

  validate() {
    this.errors = {};

    if (!this.email) {
      this.errors.email = 'Email is required.';
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(this.email)) {
      this.errors.email = 'Invalid email format.';
    }

    if (!this.password) {
      this.errors.password = 'Password is required.';
    } else if (this.password.length < 6) {
      this.errors.password = 'Password must be at least 6 characters.';
    }

    return Object.keys(this.errors).length === 0;
  },

  async register() {
    if (!this.validate()) return;

    this.loading = true;
    this.message = '';

    try {
      const response = await fetch('/api/auth/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email: this.email, password: this.password }),
      });

      if (!response.ok) {
        throw new Error('Invalid credentials');
      }

      this.message = 'Login successful!';
      // Redirect or handle success
    } catch (error) {
      this.message =
        error instanceof Error ? error.message : 'An error occurred';
    } finally {
      this.loading = false;
    }
  },
});
