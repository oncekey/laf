import React from "react";
import { HamburgerIcon } from "@chakra-ui/icons";
import { Tag } from "@chakra-ui/react";
import LeftPanel from "pages/app/[app_id]/mods/LeftPanel";

import FileTypeIcon, { FileType } from "@/components/FileTypeIcon";
import IconWrap from "@/components/IconWrap";
import Panel from "@/components/Panel";
import SectionList from "@/components/SectionList";

import useStorageStore, { TStorage } from "../../store";
import CreateBucketModal from "../CreateBucketModal";
import DeleteBucketModal from "../DeleteBucketModal";
import EditBucketModal from "../EditBucketModal";

export default function StorageListPanel() {
  const store = useStorageStore((store) => store);

  return (
    <LeftPanel>
      <Panel
        title="云存储"
        actions={[
          <CreateBucketModal key="create_modal" />,
          <IconWrap key="options" onClick={() => {}}>
            <HamburgerIcon fontSize={12} />
          </IconWrap>,
        ]}
      >
        <SectionList>
          {(store.allStorages || []).map((storage: TStorage) => {
            return (
              <SectionList.Item
                isActive={storage?.metadata?.name === store.currentStorage?.metadata?.name}
                key={storage?.metadata?.name}
                className="group"
                onClick={() => {
                  store.setCurrentStorage(storage);
                }}
              >
                <div className="relative flex-1">
                  <div className="flex justify-between">
                    <div>
                      <FileTypeIcon type={FileType.bucket} />
                      <span className="ml-2 text-base">{storage.metadata?.name}</span>
                    </div>

                    <Tag size="sm" colorScheme={"green"}>
                      {storage?.spec?.policy}
                    </Tag>
                  </div>
                </div>
                <div className="hidden group-hover:flex">
                  <EditBucketModal storage={storage} />
                  <DeleteBucketModal storage={storage} />
                </div>
              </SectionList.Item>
            );
          })}
        </SectionList>
      </Panel>
    </LeftPanel>
  );
}
