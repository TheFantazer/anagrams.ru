import { createApp } from "vue";
import { createPinia } from "pinia";

import App from "./App.vue";
import router from "./router";
import i18n from "./locales/i18n";

import "./assets/tokens.css";
import "./assets/app.css";
import "./assets/game.css";
import "./assets/pages.css";

const app = createApp(App);

app.use(createPinia());
app.use(i18n);
app.use(router);

app.mount("#app");