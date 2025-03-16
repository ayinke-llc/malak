import { Api } from "@/client/Api";
import { API_URL } from "./config";

const client = new Api({
  baseURL: API_URL,
});

export default client;
