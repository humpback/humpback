<script lang="ts" setup>
import LogoSVG from "@/assets/logo.png"

const props = withDefaults(
  defineProps<{
    direction?: "horizontal" | "vertical"
    showText?: boolean
  }>(),
  {
    direction: "vertical"
  }
)

const router = useRouter()

const logoName = "Humpback"
const logoImg = LogoSVG

function clickLogo() {
  router.push({ name: "workspace" })
}
</script>

<template>
  <div v-if="props.direction === 'horizontal'" class="logo-box" @click="clickLogo()">
    <el-avatar :size="32" :src="logoImg" class="color-light" fit="fill" shape="square" />
    <span v-if="props.showText" class="gradient-text"> {{ logoName }} </span>
  </div>
  <div v-if="props.direction === 'vertical'" class="w-100 text-align-center">
    <el-avatar :size="80" :src="logoImg" class="color-transparent" fit="fill" shape="square" />
    <div class="gradient-text mt-3"> {{ logoName }}</div>
  </div>
</template>

<style lang="scss" scoped>
.logo-box {
  height: 100%;
  width: 100%;
  display: flex;
  gap: 12px;
  align-items: center;
  justify-content: left;
  box-sizing: border-box;

  .color-light {
    background-color: #ffffff;
  }

  &:hover {
    cursor: pointer;
  }
}

.el-avatar {
  border-radius: 50%;
}

.color-transparent {
  background-color: rgba(0, 0, 0, 0);
}

.gradient-text {
  font-size: 20px;
  font-weight: bold;
  background-image: linear-gradient(45deg, #ffc700, #ff3d00);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
}
</style>
