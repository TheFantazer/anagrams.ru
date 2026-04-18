import { createI18n } from "vue-i18n";
import { useLangStore } from "@/stores/lang";

import ru from "./ru.json";
import en from "./en.json";

const i18n = createI18n({
    legacy: false,
    locale: "ru",
    fallbackLocale: "en",
    messages: { ru, en }
});

export function syncI18nWithStore(pinia) {
    const langStore = useLangStore(pinia);

    i18n.global.locale.value = langStore.lang;

    langStore.setLang(i18n.global.locale.value);
}

export default i18n;