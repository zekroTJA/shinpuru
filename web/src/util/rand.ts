export const randomNumber = (max: number, min = 0) => {
  return Math.floor(Math.floor(Math.random() * (max - min + 1)) + min);
};

export const randomFrom = <T>(array: T[]) =>
  array[randomNumber(array.length - 1)];
