import JSEncrypt from "jsencrypt"
import { RSAPublicKey } from "@/models"

export function RSAEncrypt(text: string): string {
  const encryptor = new JSEncrypt()
  encryptor.setPublicKey(RSAPublicKey)
  return encryptor.encrypt(text) as string
}

export function GenerateUUID(): string {
  let uuid = ""
  const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

  for (let i = 0; i < 32; i++) {
    const index = Math.floor(Math.random() * chars.length)
    uuid += chars[index]
    if (i === 7 || i === 11 || i === 15 || i === 19) {
      uuid += "-"
    }
  }
  return uuid
}
