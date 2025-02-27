import Alpine from 'alpinejs';
import { form } from './components';

document.addEventListener('alpine:init', () => {
  Alpine.data('form', form);
});

Alpine.start();
