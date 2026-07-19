import { dataURLFromBase64 } from "./images.ts";

export type PreparedMaskComposite = {
  sourceDataURL: string;
  maskB64: string;
};

export async function normalizeMaskForSource(
  sourceDataURL: string,
  maskB64: string,
): Promise<PreparedMaskComposite> {
  const [source, mask] = await Promise.all([
    loadImage(sourceDataURL),
    loadImage(dataURLFromBase64(maskB64)),
  ]);
  const canvas = makeCanvas(source.naturalWidth, source.naturalHeight);
  const context = requireContext(canvas);
  context.drawImage(mask, 0, 0, canvas.width, canvas.height);
  const pixels = context.getImageData(0, 0, canvas.width, canvas.height);
  let hasTransparency = false;
  for (let offset = 3; offset < pixels.data.length; offset += 4) {
    const alpha = pixels.data[offset];
    if (alpha < 255) hasTransparency = true;
  }
  if (!hasTransparency) {
    for (let offset = 0; offset < pixels.data.length; offset += 4) {
      const luma = Math.round(
        pixels.data[offset] * 0.299
        + pixels.data[offset + 1] * 0.587
        + pixels.data[offset + 2] * 0.114,
      );
      pixels.data[offset] = 0;
      pixels.data[offset + 1] = 0;
      pixels.data[offset + 2] = 0;
      pixels.data[offset + 3] = 255 - luma;
    }
    context.putImageData(pixels, 0, 0);
  }
  return { sourceDataURL, maskB64: stripDataURL(canvas.toDataURL("image/png")) };
}

export async function compositeMaskedEdit(
  prepared: PreparedMaskComposite,
  generatedB64: string,
): Promise<string> {
  const [source, generated, mask] = await Promise.all([
    loadImage(prepared.sourceDataURL),
    loadImage(dataURLFromBase64(generatedB64)),
    loadImage(dataURLFromBase64(prepared.maskB64)),
  ]);
  const width = source.naturalWidth;
  const height = source.naturalHeight;
  const maskCanvas = makeCanvas(width, height);
  const maskContext = requireContext(maskCanvas);
  maskContext.drawImage(mask, 0, 0, width, height);
  const maskPixels = maskContext.getImageData(0, 0, width, height).data;

  const sourceCanvas = makeCanvas(width, height);
  const sourceContext = requireContext(sourceCanvas);
  sourceContext.drawImage(source, 0, 0, width, height);
  const sourcePixels = sourceContext.getImageData(0, 0, width, height).data;

  const generatedCanvas = makeCanvas(width, height);
  const generatedContext = requireContext(generatedCanvas);
  generatedContext.drawImage(generated, 0, 0, width, height);
  const generatedPixels = generatedContext.getImageData(0, 0, width, height).data;

  const result = generatedContext.createImageData(width, height);
  for (let offset = 0; offset < result.data.length; offset += 4) {
    const keep = maskPixels[offset + 3];
    const edit = 255 - keep;
    for (let channel = 0; channel < 4; channel++) {
      result.data[offset + channel] = Math.round(
        (sourcePixels[offset + channel] * keep + generatedPixels[offset + channel] * edit) / 255,
      );
    }
  }
  generatedContext.putImageData(result, 0, 0);
  return stripDataURL(generatedCanvas.toDataURL("image/png"));
}

function makeCanvas(width: number, height: number): HTMLCanvasElement {
  if (!Number.isFinite(width) || !Number.isFinite(height) || width <= 0 || height <= 0) {
    throw new Error("蒙版源图尺寸无效");
  }
  const canvas = document.createElement("canvas");
  canvas.width = width;
  canvas.height = height;
  return canvas;
}

function requireContext(canvas: HTMLCanvasElement): CanvasRenderingContext2D {
  const context = canvas.getContext("2d", { willReadFrequently: true });
  if (!context) throw new Error("无法创建蒙版合成画布");
  return context;
}

function loadImage(dataURL: string): Promise<HTMLImageElement> {
  return new Promise((resolve, reject) => {
    const image = new Image();
    image.onload = () => resolve(image);
    image.onerror = () => reject(new Error("无法解码蒙版或源图片"));
    image.src = dataURL;
  });
}

function stripDataURL(value: string): string {
  const comma = value.indexOf(",");
  return comma >= 0 ? value.slice(comma + 1) : value;
}
