import { assert } from "jsr:@std/assert";

import * as ajv from "http://localhost:8080/ajv@8.12.0?bundle&external=fast-deep-equal";

function use(arg: any) {
}

Deno.test("external-bundle", async () => {
  use(ajv); // hack to prevent unused import removal
  assert((globalThis as any).ourDeepEqImported === true);
});
