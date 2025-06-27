<script setup lang="ts">
import {Page} from "@vben/common-ui";
</script>

<template>
  <Page title="Create a reaction">
    <ElCard>
      <div id="marvin-js" style="margin-top: 10px">
        <iframe
          id="ketcher-js-editor"
          src="/static/ketcher/index.html"
          width="100%"
          height="450px"
        ></iframe>
      </div>
    </ElCard>
  </Page>
</template>

<style scoped>

</style>
