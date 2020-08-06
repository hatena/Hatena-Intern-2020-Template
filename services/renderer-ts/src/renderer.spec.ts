import { render } from "./renderer";

describe("render", () => {
  it("URL の自動リンクができる", async () => {
    const src = "foo https://google.com/ bar";
    const html = await render(src);
    expect(html).toBe('foo <a href="https://google.com/">https://google.com/</a> bar');
  });
});
