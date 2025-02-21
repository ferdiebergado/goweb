interface FormValues {
  [key: string]: any;
}

interface ValidationRules {
  [key: string]: (value: any) => string | undefined;
}

interface FormHandlerOptions {
  frmId?: string;
  method?: string; // Default: POST
  validationRules?: ValidationRules;
  onSuccess?: (data: any) => void;
  onError?: (error: string | any) => void;
  onFinally?: () => void; // Called regardless of success or failure
}

export class FormHandler {
  private options: FormHandlerOptions;
  private form: HTMLFormElement;

  constructor(options: FormHandlerOptions) {
    this.options = {
      method: 'POST', // Default method
      ...options,
    };
    this.form = document.getElementById(
      this.options.frmId || 'x-form'
    ) as HTMLFormElement;
  }

  public async submit(formValues: FormValues): Promise<void> {
    try {
      // 1. Validation
      const errors = this.validate(formValues);
      if (Object.keys(errors).length > 0) {
        throw new Error(JSON.stringify(errors)); // Throw error to be caught later
      }

      // 2. API Request
      const response = await fetch(this.form.action, {
        method: this.options.method,
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(formValues),
      });

      if (!response.ok) {
        const errorData = await response.json(); // Attempt to parse error response
        throw new Error(errorData.message || response.statusText); // Use message from error response or default status text
      }

      const responseData = await response.json();

      // 3. Success Callback
      if (this.options.onSuccess) {
        this.options.onSuccess(responseData);
      }
    } catch (error) {
      // 4. Error Handling
      if (this.options.onError) {
        if (typeof error === 'string') {
          this.options.onError(error); // Handle string errors (e.g., from validation)
        } else if (error instanceof Error && error.message.startsWith('{')) {
          try {
            const parsedErrors = JSON.parse(error.message);
            this.options.onError(parsedErrors); // Handle structured validation errors
          } catch (parseError) {
            this.options.onError(error.message); // Handle parsing errors or other Error types
          }
        } else if (error instanceof Error) {
          this.options.onError(error.message); // Handle other Error types
        } else {
          this.options.onError(error);
        }
      }
    } finally {
      // 5. Finally Callback
      if (this.options.onFinally) {
        this.options.onFinally();
      }
    }
  }

  public handleSubmit(): void {
    this.form.addEventListener('submit', (event) => {
      event.preventDefault();

      const formData = new FormData(this.form);
      const formValues = Object.fromEntries(formData);
      this.submit(formValues);
    });
  }

  private validate(formValues: FormValues): { [key: string]: string } {
    const errors: { [key: string]: string } = {};
    if (this.options.validationRules) {
      for (const fieldName in this.options.validationRules) {
        const validator = this.options.validationRules[fieldName];
        const errorMessage = validator(formValues[fieldName]);
        if (errorMessage) {
          errors[fieldName] = errorMessage;
        }
      }
    }
    return errors;
  }
}
