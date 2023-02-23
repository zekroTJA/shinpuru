export enum KarmaType {
  VERY_LOW,
  LOW,
  ZERO,
  HIGH,
  VERY_HIGH,
}

export const getKarmaType = (v: number) => {
  if (v > 100) return KarmaType.VERY_HIGH;
  if (v > 0) return KarmaType.HIGH;
  if (v === 0) return KarmaType.ZERO;
  if (v < 0) return KarmaType.LOW;
  return KarmaType.VERY_LOW;
};
