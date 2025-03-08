type ValueType<T> = {
  [K in keyof T]: T[K];
};

export { ValueType };
