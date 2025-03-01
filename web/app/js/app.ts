import Alpine from 'alpinejs';
import { regForm } from './components';

document.addEventListener('alpine:init', () => {
  Alpine.data('regForm', regForm);
});

Alpine.start();
