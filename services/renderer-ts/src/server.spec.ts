import { handleRender } from "./server";
import { RenderRequest } from "../pb/renderer/renderer_pb";

describe("handleRender", () => {
  it("記法変換リクエストを受けて変換できる", async () => {
    const req = new RenderRequest();
    req.setSrc("foo https://google.com/ bar");
    const ctx = new Map<string, unknown>();
    const reply = await handleRender(req, ctx);
    expect(reply.getHtml()).toBe('foo <a href="https://google.com/">https://google.com/</a> bar');
  });
});
