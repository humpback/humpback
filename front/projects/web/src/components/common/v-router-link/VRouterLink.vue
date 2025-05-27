<script lang="ts" setup>
const props = withDefaults(
  defineProps<{
    text: string
    href: string
    underline?: boolean
    disabled?: boolean
    showTitle?: boolean
  }>(),
  {
    underline: false
  }
)

const emits = defineEmits<{
  (e: "click-route", isHref?: boolean): void
}>()

const router = useRouter()

function navigateToRoute(event: MouseEvent) {
  if (props.disabled) {
    return
  }
  if (event.ctrlKey || event.metaKey) {
    window.open(props.href, "_blank")
    emits("click-route", true)
  } else {
    router.push(props.href)
    emits("click-route", false)
  }
}

const classList = computed(() => {
  return props.underline ? ["link-style", "underline"] : ["link-style"]
})
</script>

<template>
  <a :class="classList" :href="props.href" :title="props.showTitle ? props.text : ''" @click.prevent.stop="navigateToRoute">
    <slot name="prefix-text" />
    {{ props.text }}
  </a>
</template>

<style lang="scss" scoped>
.link-style {
  display: inline-block;
  max-width: 100%;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  vertical-align: middle;
  line-height: 1.5;
  text-decoration: none;
  color: var(--el-color-primary);

  &:hover {
    opacity: 0.7;
    cursor: pointer;
  }
}

.underline:hover {
  text-decoration: underline;
  text-decoration-skip-ink: none;
  text-underline-offset: 4px;
}
</style>

>
