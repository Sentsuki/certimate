import { createContext, useContext, useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { App, Button, Card, Form, Input, Select, Skeleton } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { produce } from "immer";
import { z } from "zod";

import Show from "@/components/Show";
import Tips from "@/components/Tips";
import { type CAProviderType, CA_PROVIDERS } from "@/domain/provider";
import { SETTINGS_NAMES, type SSLProviderSettingsContent, type SettingsModel } from "@/domain/settings";
import { useAntdForm } from "@/hooks";
import { get as getSettings, save as saveSettings } from "@/repository/settings";
import { mergeCls } from "@/utils/css";
import { getErrMsg } from "@/utils/error";

const SettingsSSLProvider = () => {
  const { t } = useTranslation();

  const { message, notification } = App.useApp();

  const [formInst] = Form.useForm<{ provider?: string }>();
  const [formPending, setFormPending] = useState(false);

  const [settings, setSettings] = useState<SettingsModel<SSLProviderSettingsContent>>();
  const [loading, setLoading] = useState(true);
  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);

      const settings = await getSettings<SSLProviderSettingsContent>(SETTINGS_NAMES.SSL_PROVIDER);
      setSettings(settings);
      setProviderValue(settings.content?.provider || CA_PROVIDERS.LETSENCRYPT);

      setLoading(false);
    };

    fetchData();
  }, []);

  const providers = [
    [CA_PROVIDERS.LETSENCRYPT, "provider.letsencrypt", "letsencrypt.org", "/imgs/providers/letsencrypt.svg"],
    [CA_PROVIDERS.LETSENCRYPTSTAGING, "provider.letsencryptstaging", "letsencrypt.org", "/imgs/providers/letsencrypt.svg"],
    [CA_PROVIDERS.ACTALISSSL, "provider.actalisssl", "actalis.com", "/imgs/providers/actalisssl.png"],
    [CA_PROVIDERS.GLOBALSIGNATLAS, "provider.globalsignatlas", "atlas.globalsign.com", "/imgs/providers/globalsignatlas.png"],
    [CA_PROVIDERS.GOOGLETRUSTSERVICES, "provider.googletrustservices", "pki.goog", "/imgs/providers/google.svg"],
    [CA_PROVIDERS.SECTIGO, "provider.sectigo", "sectigo.com", "/imgs/providers/sectigo.svg"],
    [CA_PROVIDERS.SSLCOM, "provider.sslcom", "ssl.com", "/imgs/providers/sslcom.svg"],
    [CA_PROVIDERS.ZEROSSL, "provider.zerossl", "zerossl.com", "/imgs/providers/zerossl.svg"],
    [CA_PROVIDERS.ACMECA, "provider.acmeca", "ACME v2 (RFC 8555)", "/imgs/providers/acmeca.svg"],
  ].map(([value, name, description, icon]) => {
    return {
      value: value as CAProviderType,
      name: t(name),
      description,
      icon,
    };
  });
  const [providerValue, setProviderValue] = useState<CAProviderType>(CA_PROVIDERS.LETSENCRYPT);
  const providerFormEl = useMemo(() => {
    switch (providerValue) {
      case CA_PROVIDERS.LETSENCRYPT:
        return <InternalSettingsFormProviderLetsEncrypt />;
      case CA_PROVIDERS.LETSENCRYPTSTAGING:
        return <InternalSettingsFormProviderLetsEncryptStaging />;
      case CA_PROVIDERS.ACTALISSSL:
        return <InternalSettingsFormProviderActalisSSL />;
      case CA_PROVIDERS.GLOBALSIGNATLAS:
        return <InternalSettingsFormProviderGlobalSignAtlas />;
      case CA_PROVIDERS.GOOGLETRUSTSERVICES:
        return <InternalSettingsFormProviderGoogleTrustServices />;
      case CA_PROVIDERS.SECTIGO:
        return <InternalSettingsFormProviderSectigo />;
      case CA_PROVIDERS.SSLCOM:
        return <InternalSettingsFormProviderSSLCom />;
      case CA_PROVIDERS.ZEROSSL:
        return <InternalSettingsFormProviderZeroSSL />;
      case CA_PROVIDERS.ACMECA:
        return <InternalSettingsFormProviderACMECA />;
    }
  }, [providerValue]);

  const updateContextSettings = async (settings: MaybeModelRecordWithId<SettingsModel<SSLProviderSettingsContent>>) => {
    setFormPending(true);

    try {
      const resp = await saveSettings(settings);
      setSettings(resp);
      setProviderValue(resp.content?.provider);

      message.success(t("common.text.operation_succeeded"));
    } catch (err) {
      notification.error({ message: t("common.text.request_error"), description: getErrMsg(err) });
    } finally {
      setFormPending(false);
    }
  };

  return (
    <InternalSettingsContext.Provider
      value={{
        pending: formPending,
        settings: settings!,
        updateSettings: updateContextSettings,
      }}
    >
      <h2>{t("settings.sslprovider.ca.title")}</h2>
      <Show when={!loading} fallback={<Skeleton active />}>
        <Form form={formInst} disabled={formPending} layout="vertical" initialValues={{ provider: providerValue }}>
          <Form.Item>
            <Tips message={<span dangerouslySetInnerHTML={{ __html: t("settings.sslprovider.ca.tips") }}></span>} />
          </Form.Item>

          <Form.Item name="provider" label={t("settings.sslprovider.form.provider.label")} extra={t("settings.sslprovider.form.provider.help")}>
            <div className="flex w-full flex-wrap items-center gap-4">
              {providers.map((provider) => (
                <Card
                  key={provider.value}
                  className={mergeCls("relative overflow-hidden", { ["border-primary"]: providerValue === provider.value })}
                  style={{ width: 280 }}
                  styles={{
                    body: { padding: 0 },
                  }}
                  hoverable
                  onClick={() => setProviderValue(provider.value)}
                >
                  <div className="relative z-1 px-3 py-4">
                    <div className="flex items-center justify-between gap-3">
                      <div>
                        <img src={provider.icon} className="size-8" />
                      </div>
                      <div className="flex-1 overflow-hidden">
                        <div className="truncate">{provider.name}</div>
                        <div className="mt-1 truncate text-xs">{provider.description}</div>
                      </div>
                    </div>
                  </div>
                  {providerValue === provider.value && <div className="absolute top-0 left-0 size-full bg-primary opacity-20"></div>}
                </Card>
              ))}
            </div>
          </Form.Item>
        </Form>

        <div className="md:max-w-160">{providerFormEl}</div>
      </Show>
    </InternalSettingsContext.Provider>
  );
};

const InternalSettingsContext = createContext(
  {} as {
    pending: boolean;
    settings: SettingsModel<SSLProviderSettingsContent>;
    updateSettings: (settings: MaybeModelRecordWithId<SettingsModel<SSLProviderSettingsContent>>) => Promise<void>;
  }
);

const InternalSharedForm = ({ children, provider }: { children?: React.ReactNode; provider: CAProviderType }) => {
  const { t } = useTranslation();

  const { pending, settings, updateSettings } = useContext(InternalSettingsContext);

  const { form: formInst, formProps } = useAntdForm<NonNullable<unknown>>({
    initialValues: settings?.content?.config?.[provider],
    onSubmit: async (values) => {
      const newSettings = produce(settings, (draft) => {
        draft.content ??= {} as SSLProviderSettingsContent;
        draft.content.provider = provider;

        draft.content.config ??= {} as SSLProviderSettingsContent["config"];
        draft.content.config[provider] = values;
      });
      await updateSettings(newSettings);

      setFormChanged(false);
    },
  });

  const [formChanged, setFormChanged] = useState(false);
  useEffect(() => {
    setFormChanged(provider !== settings?.content?.provider);
  }, [provider, settings?.content?.provider]);

  const handleFormChange = () => {
    setFormChanged(true);
  };

  return (
    <Form {...formProps} form={formInst} disabled={pending} layout="vertical" onValuesChange={handleFormChange}>
      {children}

      <Form.Item>
        <Button type="primary" htmlType="submit" disabled={!formChanged} loading={pending}>
          {t("common.button.save")}
        </Button>
      </Form.Item>
    </Form>
  );
};

const InternalSharedFormEabFields = ({ i18nKey }: { i18nKey: string }) => {
  const { t, i18n } = useTranslation();

  const hasGuide = i18n.exists(`access.form.${i18nKey}_eab.guide`);

  const formSchema = z.object({
    endpoint: z.url(t("common.errmsg.url_invalid")),
    eabKid: z
      .string(t("access.form.shared_acme_eab_kid.label"))
      .min(1, t("access.form.shared_acme_eab_kid.placeholder"))
      .max(256, t("common.errmsg.string_max", { max: 256 })),
    eabHmacKey: z
      .string(t("access.form.shared_acme_eab_hmac_key.label"))
      .min(1, t("access.form.shared_acme_eab_hmac_key.placeholder"))
      .max(256, t("common.errmsg.string_max", { max: 256 })),
  });
  const formRule = createSchemaFieldRule(formSchema);

  return (
    <>
      <Form.Item name="eabKid" label={t("access.form.shared_acme_eab_kid.label")} rules={[formRule]}>
        <Input autoComplete="new-password" placeholder={t("access.form.shared_acme_eab_kid.placeholder")} />
      </Form.Item>

      <Form.Item name="eabHmacKey" label={t("access.form.shared_acme_eab_hmac_key.label")} rules={[formRule]}>
        <Input.Password autoComplete="new-password" placeholder={t("access.form.shared_acme_eab_hmac_key.placeholder")} />
      </Form.Item>

      <Form.Item hidden={!hasGuide}>
        <Tips message={<span dangerouslySetInnerHTML={{ __html: t(`access.form.${i18nKey}_eab.guide`) }}></span>} />
      </Form.Item>
    </>
  );
};

const InternalSettingsFormProviderLetsEncrypt = () => {
  return <InternalSharedForm provider={CA_PROVIDERS.LETSENCRYPT} />;
};

const InternalSettingsFormProviderLetsEncryptStaging = () => {
  const { t } = useTranslation();

  return (
    <InternalSharedForm provider={CA_PROVIDERS.LETSENCRYPTSTAGING}>
      <Form.Item>
        <Tips message={<span dangerouslySetInnerHTML={{ __html: t("settings.sslprovider.form.letsencryptstaging_alert") }}></span>} />
      </Form.Item>
    </InternalSharedForm>
  );
};

const InternalSettingsFormProviderActalisSSL = () => {
  return (
    <InternalSharedForm provider={CA_PROVIDERS.ACTALISSSL}>
      <InternalSharedFormEabFields i18nKey="actalisssl" />
    </InternalSharedForm>
  );
};

const InternalSettingsFormProviderGlobalSignAtlas = () => {
  return (
    <InternalSharedForm provider={CA_PROVIDERS.GLOBALSIGNATLAS}>
      <InternalSharedFormEabFields i18nKey="globalsignatlas" />
    </InternalSharedForm>
  );
};

const InternalSettingsFormProviderGoogleTrustServices = () => {
  return (
    <InternalSharedForm provider={CA_PROVIDERS.GOOGLETRUSTSERVICES}>
      <InternalSharedFormEabFields i18nKey="googletrustservices" />
    </InternalSharedForm>
  );
};

const InternalSettingsFormProviderSectigo = () => {
  const { t } = useTranslation();

  const formSchema = z.object({
    validationType: z.string().nonempty(t("access.form.sectigo_validation_type.placeholder")),
  });
  const formRule = createSchemaFieldRule(formSchema);

  return (
    <InternalSharedForm provider={CA_PROVIDERS.SECTIGO}>
      <Form.Item name="validationType" initialValue="dv" label={t("access.form.sectigo_validation_type.label")} rules={[formRule]}>
        <Select
          options={["dv", "ov", "ev"].map((s) => ({
            key: s,
            label: t(`access.form.sectigo_validation_type.option.${s}.label`),
            value: s,
          }))}
          placeholder={t("access.form.sectigo_validation_type.placeholder")}
        />
      </Form.Item>

      <InternalSharedFormEabFields i18nKey="sectigo" />
    </InternalSharedForm>
  );
};

const InternalSettingsFormProviderSSLCom = () => {
  return (
    <InternalSharedForm provider={CA_PROVIDERS.SSLCOM}>
      <InternalSharedFormEabFields i18nKey="sslcom" />
    </InternalSharedForm>
  );
};

const InternalSettingsFormProviderZeroSSL = () => {
  return (
    <InternalSharedForm provider={CA_PROVIDERS.ZEROSSL}>
      <InternalSharedFormEabFields i18nKey="zerossl" />
    </InternalSharedForm>
  );
};

const InternalSettingsFormProviderACMECA = () => {
  const { t } = useTranslation();

  const formSchema = z.object({
    endpoint: z.url(t("common.errmsg.url_invalid")),
    eabKid: z
      .string(t("access.form.acmeca_eab_kid.placeholder"))
      .max(256, t("common.errmsg.string_max", { max: 256 }))
      .nullish(),
    eabHmacKey: z
      .string(t("access.form.acmeca_eab_hmac_key.placeholder"))
      .max(256, t("common.errmsg.string_max", { max: 256 }))
      .nullish(),
  });
  const formRule = createSchemaFieldRule(formSchema);

  return (
    <InternalSharedForm provider={CA_PROVIDERS.ACMECA}>
      <Form.Item
        name="endpoint"
        label={t("access.form.acmeca_endpoint.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.acmeca_endpoint.tooltip") }}></span>}
      >
        <Input placeholder={t("access.form.acmeca_endpoint.placeholder")} />
      </Form.Item>

      <Form.Item
        name="eabKid"
        label={t("access.form.acmeca_eab_kid.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.acmeca_eab_kid.tooltip") }}></span>}
      >
        <Input autoComplete="new-password" placeholder={t("access.form.acmeca_eab_kid.placeholder")} />
      </Form.Item>

      <Form.Item
        name="eabHmacKey"
        label={t("access.form.acmeca_eab_hmac_key.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.acmeca_eab_hmac_key.tooltip") }}></span>}
      >
        <Input.Password autoComplete="new-password" placeholder={t("access.form.acmeca_eab_hmac_key.placeholder")} />
      </Form.Item>
    </InternalSharedForm>
  );
};

export default SettingsSSLProvider;
