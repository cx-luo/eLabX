import { Ketcher } from 'ketcher-core';

export enum ChemicalMimeType {
  Mol = 'chemical/x-mdl-molfile',
  Rxn = 'chemical/x-mdl-rxnfile',
  DaylightSmiles = 'chemical/x-daylight-smiles',
  ExtendedSmiles = 'chemical/x-chemaxon-cxsmiles',
  DaylightSmarts = 'chemical/x-daylight-smarts',
  InChI = 'chemical/x-inchi',
  InChIAuxInfo = 'chemical/x-inchi-aux',
  InChIKey = 'chemical/x-inchi-key',
  CDX = 'chemical/x-cdx',
  CDXML = 'chemical/x-cdxml',
  CML = 'chemical/x-cml',
  KET = 'chemical/x-indigo-ket',
  UNKNOWN = 'chemical/x-unknown',
  SDF = 'chemical/x-sdf',
  FASTA = 'chemical/x-fasta',
  SEQUENCE = 'chemical/x-sequence',
  PeptideSequenceThreeLetter = 'chemical/x-peptide-sequence-3-letter',
  RNA = 'chemical/x-rna-sequence',
  DNA = 'chemical/x-dna-sequence',
  PEPTIDE = 'chemical/x-peptide-sequence',
  IDT = 'chemical/x-idt',
  HELM = 'chemical/x-helm',
  RDF = 'chemical/x-rdf',
}

// export declare class Ketcher {
//   #private;
//   logging: LogSettings;
//   structService: StructService;
//   _indigo: Indigo;
//   changeEvent: Subscription;
//   get editor(): Editor;
//   get eventBus(): EventEmitter;
//   constructor(editor: Editor, structService: StructService, formatterFactory: FormatterFactory);
//   get formatterFactory(): FormatterFactory;
//   get indigo(): Indigo;
//   get settings(): {};
//   setSettings(settings: Record<string, string>): any;
//   getSmiles(isExtended?: boolean): Promise<string>;
//   getMolfile(molfileFormat?: MolfileFormat): Promise<string>;
//   getIdt(): Promise<string>;
//   getRxn(molfileFormat?: MolfileFormat): Promise<string>;
//   getKet(): Promise<string>;
//   getFasta(): Promise<string>;
//   getSequence(): Promise<string>;
//   getSmarts(): Promise<string>;
//   getCml(): Promise<string>;
//   getSdf(molfileFormat?: MolfileFormat): Promise<string>;
//   getRdf(molfileFormat?: MolfileFormat): Promise<string>;
//   getCDXml(): Promise<string>;
//   getCDX(): Promise<string>;
//   getInchi(withAuxInfo?: boolean): Promise<string>;
//   getInChIKey(): Promise<string>;
//   containsReaction(): boolean;
//   isQueryStructureSelected(): boolean;
//   setMolecule(structStr: string, options?: SetMoleculeOptions): Promise<void | undefined>;
//   setHelm(helmStr: string): Promise<void | undefined>;
//   addFragment(structStr: string, options?: SetMoleculeOptions): Promise<void | undefined>;
//   layout(): Promise<void>;
//   calculate(options?: CalculateData): Promise<CalculateResult>;
//   /**
//    * @param {number} value - in a range [ZoomTool.instance.MINZOOMSCALE, ZoomTool.instance.MAXZOOMSCALE]
//    */
//   setZoom(value: number): void;
//   setMode(mode: SupportedModes): void;
//   exportImage(format: SupportedImageFormats, params?: ExportImageParams): void;
//   recognize(image: Blob, version?: string): Promise<Struct>;
//   generateImage(data: string, options?: GenerateImageOptions): Promise<Blob>;
//   reinitializeIndigo(structService: StructService): void;
//   sendCustomAction(name: string): void;
//   updateMonomersLibrary(rawMonomersData: string | JSON): void;
// }

// types/ketcher.d.ts
declare module 'ketcher-core' {
  interface Ketcher {
    /**
     * 重载方法：支持不同参数返回不同结构
     */
    getInChIKey(options?: { detailed?: false }): Promise<string>;

    getInChIKey(options: { detailed: true }): Promise<ApiResponse>;
  }
}

export interface ApiResponse {
  format: string;
  original_format: string;
  struct: string;
}

/**
 * 获取 Ketcher 实例
 * @returns Ketcher | null - 返回实例或 null 表示失败
 */
export function getKetcher(): Ketcher | null {
  const ketcherFrame = document.getElementById('ketcher-js-editor') as HTMLIFrameElement | null;

  if (ketcherFrame && 'contentDocument' in ketcherFrame) {
    const contentWindow = ketcherFrame.contentWindow;

    if (contentWindow && 'ketcher' in contentWindow) {
      return contentWindow.ketcher as Ketcher;
    }
  }

  return null; // 如果任何一步失败，返回 null
}

export async function exportInChiKey(ketcher: Ketcher): Promise<string | null> {
  try {
    const response: ApiResponse = await ketcher.getInChIKey({ detailed: true });
    return response.struct;
  } catch (e: any) {
    console.error('getInChIKey error', e);
    return null; // 或者根据业务需求返回其他默认值
  }
}

export function generateImgUrl(ketcher: Ketcher, smiles: string) {
  let imgUrl: string | null = null;
  if (smiles) {
    ketcher
      .generateImage(smiles, {
        outputFormat: 'svg', // 生成图片的类型，可以是"svg"或"png"
        backgroundColor: '255, 255, 255', // 背景颜色
      })
      .then((res: Blob) => {
        imgUrl = window.URL.createObjectURL(res); // res是blob类型，用该方法转为url后可以在用img展示
      });
  } else {
    imgUrl = null;
  }
  return imgUrl;
}
