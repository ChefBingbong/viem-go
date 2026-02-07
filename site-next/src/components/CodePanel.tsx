"use client";

import * as React from "react";
import { Highlight, themes } from "prism-react-renderer";
import {
  Tab,
  TabGroup,
  TabList,
  TabPanel,
  TabPanels,
} from "@headlessui/react";
import CopyButton from "./CopyButton";

const languageNames: Record<string, string> = {
  js: "JavaScript",
  jsx: "JavaScript",
  ts: "TypeScript",
  tsx: "TypeScript",
  javascript: "JavaScript",
  typescript: "TypeScript",
  go: "Go",
  bash: "Terminal",
  shell: "Terminal",
  json: "JSON",
  yaml: "YAML",
  css: "CSS",
  html: "HTML",
};

interface CodeTab {
  code: string;
  language?: string;
  title?: string;
  showLineNumbers?: boolean;
}

interface CodeGroupProps {
  tabs: CodeTab[];
  title?: string;
}

export function CodeGroup({ tabs: tabsInput, title }: CodeGroupProps) {
  const [selectedIndex, setSelectedIndex] = React.useState(0);

  const tabs = tabsInput.map((tab) => ({
    title:
      tab.title ||
      languageNames[tab.language || ""] ||
      tab.language ||
      "Code",
    language: tab.language || "typescript",
    code: tab.code || "",
    showLineNumbers: tab.showLineNumbers ?? true,
  }));

  if (tabs.length === 0) return null;

  if (tabs.length === 1) {
    const tab = tabs[0]!;
    const codeStr = (tab.code || "").trim();
    return (
      <div className="my-6 rounded-lg overflow-hidden border border-accent/20 bg-gray-6/80">
        <div className="flex items-center justify-between px-3 h-11 bg-dark-deep/60">
          <span className="font-mono text-[0.8125rem] font-medium text-white leading-none">
            {title || tab.title}
          </span>
          <CopyButton text={codeStr} />
        </div>
        <Highlight
          theme={themes.nightOwl}
          code={codeStr}
          language={tab.language}
        >
          {({ tokens, getLineProps, getTokenProps }) => (
            <pre className="m-0 py-1.5 px-3 overflow-auto text-[0.8125rem] leading-relaxed bg-transparent">
              {tokens.map((line, i) => (
                <div
                  key={i}
                  {...getLineProps({ line })}
                  className="table-row m-0"
                >
                  {tab.showLineNumbers && (
                    <span className="table-cell pr-3 text-right text-gray-4 select-none min-w-6">
                      {i + 1}
                    </span>
                  )}
                  <span className="table-cell">
                    {line.map((token, key) => (
                      <span key={key} {...getTokenProps({ token })} />
                    ))}
                  </span>
                </div>
              ))}
            </pre>
          )}
        </Highlight>
      </div>
    );
  }

  return (
    <div className="my-6 rounded-lg overflow-hidden border border-accent/20 bg-gray-6/80">
      <TabGroup selectedIndex={selectedIndex} onChange={setSelectedIndex}>
        <div className="flex items-center justify-between h-11 bg-dark-deep/90">
          <TabList className="flex h-full items-stretch">
            {tabs.map((tab, index) => (
              <Tab
                key={index}
                className={`flex items-center justify-center px-3.5 text-[0.8125rem] font-medium cursor-pointer transition-all duration-150 h-11 border-b-2 outline-none ${
                  selectedIndex === index
                    ? "text-white bg-white/5 border-accent"
                    : "text-gray-3 bg-transparent border-transparent hover:text-white"
                }`}
              >
                {tab.title}
              </Tab>
            ))}
          </TabList>
          <div className="flex items-center pr-2">
            <CopyButton text={(tabs[selectedIndex]?.code || "").trim()} />
          </div>
        </div>

        <TabPanels className="mt-0">
          {tabs.map((tab, index) => {
            const codeStr = (tab.code || "").trim();
            return (
              <TabPanel key={index}>
                <Highlight
                  theme={themes.nightOwl}
                  code={codeStr}
                  language={tab.language}
                >
                  {({ tokens, getLineProps, getTokenProps }) => (
                    <pre className="m-0 py-1.5 px-3 overflow-auto text-[0.8125rem] leading-relaxed bg-transparent">
                      {tokens.map((line, i) => (
                        <div
                          key={i}
                          {...getLineProps({ line })}
                          className="table-row m-0"
                        >
                          {tab.showLineNumbers && (
                            <span className="table-cell pr-3 text-right text-gray-4 select-none min-w-6">
                              {i + 1}
                            </span>
                          )}
                          <span className="table-cell">
                            {line.map((token, key) => (
                              <span key={key} {...getTokenProps({ token })} />
                            ))}
                          </span>
                        </div>
                      ))}
                    </pre>
                  )}
                </Highlight>
              </TabPanel>
            );
          })}
        </TabPanels>
      </TabGroup>
    </div>
  );
}
