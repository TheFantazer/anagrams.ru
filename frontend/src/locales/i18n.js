import { createI18n } from "vue-i18n";

import ru from "./ru.json";
import en from "./en.json";

// Get saved language preference or default to English
const savedLanguage = localStorage.getItem('anagram_ui_lang') || 'en';

const i18n = createI18n({
  legacy: false,
  locale: savedLanguage,
  fallbackLocale: "en",
  messages: { ru, en }
});

export default i18n;
