<script setup>
</script>

<template>
  <div class="page">
    <div class="shell">
      <header class="page-head" style="margin-bottom: 24px">
        <h1 class="page-title-display">{{ $t('privacy.title') }}</h1>
      </header>

      <div class="card" style="padding: 0; overflow: hidden">
        <iframe
          src="/privacy_policy_anagrams_fixed.pdf"
          width="100%"
          height="800px"
          style="border: none; display: block"
          title="Privacy Policy"
        ></iframe>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* All styles are in pages.css and app.css */
</style>
