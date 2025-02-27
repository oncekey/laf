import React, { useState } from "react";

import FunctionLayout from "@/components/Layout/Function";
import { Pages, SiderBarWidth } from "@/constants/index";

import SiderBar from "./mods/SiderBar";
import DatabasePage from "./database";
import FunctionPage from "./functions";
import LogsPage from "./logs";
import StoragePage from "./storages";

function AppDetail() {
  const [pageId, setPageId] = useState("function");
  return (
    <>
      <SiderBar pageId={pageId} setPageId={(pageId: string) => setPageId(pageId)} />
      <div className="h-full" style={{ marginLeft: SiderBarWidth }}>
        {[
          {
            pageId: Pages.function,
            component: FunctionPage,
          },
          {
            pageId: Pages.database,
            component: DatabasePage,
          },
          {
            pageId: Pages.storage,
            component: StoragePage,
          },
          {
            pageId: Pages.logs,
            component: LogsPage,
          },
        ].map((item) => (
          <div key={item.pageId} className={pageId === item.pageId ? "block h-full" : "hidden"}>
            {item.pageId === pageId ? <item.component /> : null}
          </div>
        ))}
      </div>
    </>
  );
}

AppDetail.getLayout = (page: React.ReactElement) => {
  return <FunctionLayout>{page}</FunctionLayout>;
};

export default AppDetail;
