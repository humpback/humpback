<script lang="ts" setup>
const props = withDefaults(defineProps<{ size?: number }>(), { size: 18 })
</script>

<template>
  <div :style="{ width: `${props.size}px`, height: `${props.size}px` }" class="loader" />
</template>

<style lang="scss" scoped>
/* HTML: <div class="loader"></div> */
.loader {
  width: 30px;
  aspect-ratio: 1;
  display: grid;
  -webkit-mask: conic-gradient(from 15deg, #0000, #000);
  animation: l26 1.2s infinite steps(12);
}

.loader,
.loader:before,
.loader:after {
  background:
    radial-gradient(closest-side at 50% 12.5%, var(--el-color-info) 96%, #0000) 50% 0/20% 80% repeat-y,
    radial-gradient(closest-side at 12.5% 50%, var(--el-color-info) 96%, #0000) 0 50%/80% 20% repeat-x;
}

.loader:before,
.loader:after {
  content: "";
  grid-area: 1/1;
  transform: rotate(30deg);
}

.loader:after {
  transform: rotate(60deg);
}

@keyframes l26 {
  100% {
    transform: rotate(1turn);
  }
}
</style>
