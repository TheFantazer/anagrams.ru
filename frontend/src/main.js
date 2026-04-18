import { createApp } from "vue";
import { createPinia } from "pinia";

import App from "./App.vue";
import router from "./router";
import i18n from "./locales/i18n";
import { useUserStore } from "./stores/userStore";

import "./assets/tokens.css";
import "./assets/app.css";
import "./assets/game.css";
import "./assets/pages.css";

const app = createApp(App);
const pinia = createPinia();

app.use(pinia);
app.use(i18n);
app.use(router);

const userStore = useUserStore();
i18n.global.locale.value = userStore.soloLang;

import { watch } from "vue";
watch(() => userStore.soloLang, (newLang) => {
  i18n.global.locale.value = newLang;
});

app.mount("#app");