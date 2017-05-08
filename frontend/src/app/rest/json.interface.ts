export interface JsonInterface {
  id: string;
  json(): string;
}

export interface JsonConstructor {
  new (json: Object): JsonInterface;
}

export function createModel(ctor: JsonConstructor, json: Object): JsonInterface {
  return new ctor(json);
}

export function Model(ctor: JsonConstructor) {
    return function (target: any, propertyKey: string, descriptor: PropertyDescriptor) {
        let originalMethod = descriptor.value;

        descriptor.value = function(...args: any[]) {
          let result = originalMethod.apply(this, args, ctor);
          return result;
        };

    return descriptor;

    };
}