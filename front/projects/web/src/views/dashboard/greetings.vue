<script lang="ts" setup>
const props = defineProps<{ ownerGroups?: number; ownerServices?: number; isLoading?: boolean }>()

const { t } = useI18n()
const userStore = useUserStore()

const greetings = computed(() => {
  const currentHour = new Date().getHours()
  if (currentHour >= 6 && currentHour < 9) {
    return {
      i18nLabel: "tips.breakfastTips",
      icon: IconNotoSalutingFace
    }
  } else if (currentHour >= 9 && currentHour < 12) {
    return {
      i18nLabel: "tips.morningTips",
      icon: IconNotoGrinningFace
    }
  } else if (currentHour >= 12 && currentHour < 14) {
    return {
      i18nLabel: "tips.middayTips",
      icon: IconNotoWinkingFace
    }
  } else if (currentHour >= 14 && currentHour < 18) {
    return {
      i18nLabel: "tips.afternoonTips",
      icon: IconNotoMeltingFace
    }
  } else if (currentHour >= 18 && currentHour < 20) {
    return {
      i18nLabel: "tips.eveningTips",
      icon: IconNotoSmilingFaceWithSunglasses
    }
  } else if (currentHour >= 20 && currentHour < 23) {
    return {
      i18nLabel: "tips.eveningAfterTips",
      icon: IconNotoYawningFace
    }
  } else {
    return {
      i18nLabel: "tips.morningBeforeTips",
      icon: IconNotoSleepingFace
    }
  }
})
</script>

<template>
  <v-card>
    <div class="greeting-box">
      <div>
        <el-icon :size="60">
          <component :is="greetings.icon" />
        </el-icon>
      </div>
      <div class="greeting-title">
        <div class="greeting-title-left">
          <div> {{ t(greetings.i18nLabel, { name: userStore.userInfo.username }) }}</div>
          <el-text style="font-weight: normal">{{ t("tips.thankUseTips") }}</el-text>
        </div>
        <div class="greeting-title-right">
          <div class="greeting-title-right-item">
            <div>{{ t("label.groups") }}</div>
            <div class="mt-3 f-semiBold">
              <span v-if="!props.isLoading">{{ props.ownerGroups || 0 }}</span>
              <v-loading v-else />
            </div>
          </div>
          <div>
            <div>{{ t("label.services") }}</div>
            <div class="mt-3 f-semiBold">
              <span v-if="!props.isLoading">{{ props.ownerServices || 0 }}</span>
              <v-loading v-else />
            </div>
          </div>
        </div>
      </div>
    </div>
  </v-card>
</template>

<style lang="scss" scoped>
.greeting-box {
  display: flex;
  align-items: start;
  gap: 12px;

  .greeting-title {
    flex: 1;
    display: flex;
    align-items: start;
    gap: 20px;
    flex-wrap: wrap;
    padding-top: 8px;

    .greeting-title-left {
      flex: 1;
      font-size: 20px;
      font-weight: bold;
    }

    .greeting-title-right {
      display: flex;
      align-items: start;
      gap: 20px;
      padding-right: 50px;

      .greeting-title-right-item {
        min-width: 80px;
      }
    }
  }
}
</style>
