export async function init() {
  const config = await commonService.config()
  initRSA(config.EncryptionKey)
  initRule(config.RuleLengthLimit, config.RuleFormat)
}
