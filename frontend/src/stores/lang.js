import { defineStore } from "pinia";
import i18n from "@/locales/i18n";

export const useLangStore = defineStore("lang", {
    state: () => ({
        lang: localStorage.getItem("lang") || "ru"
    }),

    actions: {
        setLang(value) {
            this.lang = value;
            localStorage.setItem("lang", value);

            i18n.global.locale.value = value;
        },

        toggle() {
            this.setLang(this.lang === "ru" ? "en" : "ru");
        }
    }
});