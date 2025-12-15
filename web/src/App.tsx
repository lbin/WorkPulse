import React, { useEffect, useMemo, useState } from "react";
import { useRoutes } from "react-router-dom";
import { ConfigProvider } from "antd";
import enUS from "antd/locale/en_US";
import zhCN from "antd/locale/zh_CN";

import { routes } from "./routes";
import i18n from "./i18n";

export default function App() {
  const element = useRoutes(routes);
  const [lang, setLang] = useState(i18n.language?.startsWith("en") ? "en" : "zh");

  useEffect(() => {
    const onChange = (lng: string) => setLang(lng?.startsWith("en") ? "en" : "zh");
    i18n.on("languageChanged", onChange);
    return () => { i18n.off("languageChanged", onChange); };
  }, []);

  const antdLocale = useMemo(() => (lang === "en" ? enUS : zhCN), [lang]);

  return (
    <ConfigProvider locale={antdLocale}>
      {element}
    </ConfigProvider>
  );
}
