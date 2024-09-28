const GOOGLE_CLIENT_ID = process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID;
const MALAK_TERMS_CONDITION_LINK =
	process.env.NEXT_PUBLIC_MALAK_TERMS_CONDITION_LINK ||
	"https://ayinke.ventures/malak/terms";
const MALAK_PRIVACY_POLICY_LINK =
	process.env.NEXT_PUBLIC_MALAK_PRIVACY_POLICY_LINK ||
	"https://ayinke.ventures/malak/privacy";

export {
	GOOGLE_CLIENT_ID,
	MALAK_TERMS_CONDITION_LINK,
	MALAK_PRIVACY_POLICY_LINK,
};
