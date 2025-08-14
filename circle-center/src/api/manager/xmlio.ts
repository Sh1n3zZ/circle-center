import { post, request } from "../client";
import type {
  ParseXmlRequest,
  ParseXmlResponse,
  ConfirmImportRequest,
  ConfirmImportResponse,
} from "./types";

export const xmlioApi = {
  async parse(body: ParseXmlRequest): Promise<ParseXmlResponse> {
    const res = await post<ParseXmlResponse>("/manager/icons/parse", body);
    return res.data;
  },

  // Alternative: form data version for large files (not used by default)
  async parseForm(appfilter?: File, appmap?: File, theme?: File): Promise<ParseXmlResponse> {
    const form = new FormData();
    if (appfilter) form.append("appfilter", await appfilter.text());
    if (appmap) form.append("appmap", await appmap.text());
    if (theme) form.append("theme", await theme.text());
    const res = await request({ url: "/manager/icons/parse", method: "POST", data: form, asFormData: true });
    return res.data as unknown as ParseXmlResponse;
  },

  // Step 2: confirm import
  async confirmImport(body: ConfirmImportRequest): Promise<ConfirmImportResponse> {
    const res = await post<ConfirmImportResponse>("/manager/icons/import", body);
    return res.data;
  },
};


