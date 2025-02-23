import { defaultProps } from "@blocknote/core";
import { createReactBlockSpec } from "@blocknote/react";
import { Menu } from "@mantine/core";
import { RiDashboard2Line } from "@remixicon/react";
import "./styles.css";

export const dashboardItems = [
  { id: 1, title: "Revenue Overview", value: "revenue" },
  { id: 2, title: "User Analytics", value: "users" },
  { id: 3, title: "Sales Report", value: "sales" },
  { id: 4, title: "Performance Metrics", value: "performance" },
  { id: 5, title: "Customer Feedback", value: "feedback" },
  { id: 6, title: "Inventory Status", value: "inventory" },
  { id: 7, title: "Marketing ROI", value: "marketing" },
  { id: 8, title: "Team Progress", value: "team" },
  { id: 9, title: "Project Timeline", value: "projects" },
  { id: 10, title: "Financial Summary", value: "finance" },
] as const;

export const Dashboard = createReactBlockSpec(
  {
    type: "dashboard",
    propSchema: {
      textAlignment: defaultProps.textAlignment,
      textColor: defaultProps.textColor,
      selectedItem: {
        default: "",
        values: ["", ...dashboardItems.map(item => item.value)],
      },
    },
    content: "inline",
  },
  {
    render: (props) => {
      const selectedItem = dashboardItems.find(
        (item) => item.value === props.block.props.selectedItem
      );

      return (
        <div className={"dashboard"} data-dashboard-type={props.block.props.selectedItem || "none"}>
          <Menu withinPortal={false}>
            <Menu.Target>
              <div className={"dashboard-icon-wrapper"} contentEditable={false}>
                <RiDashboard2Line
                  className={"dashboard-icon"}
                  size={32}
                />
                <span className="dashboard-title">
                  {selectedItem ? selectedItem.title : "Select Dashboard Item"}
                </span>
              </div>
            </Menu.Target>
            <Menu.Dropdown>
              <Menu.Label>Select Dashboard Item</Menu.Label>
              <Menu.Divider />
              {dashboardItems.map((item) => (
                <Menu.Item
                  key={item.value}
                  onClick={() =>
                    props.editor.updateBlock(props.block, {
                      type: "dashboard",
                      props: { selectedItem: item.value },
                    })
                  }>
                  {item.title}
                </Menu.Item>
              ))}
            </Menu.Dropdown>
          </Menu>
          {selectedItem ? (
            <div className={"inline-content"} ref={props.contentRef} />
          ) : (
            <div className={"dashboard-placeholder"}>Click the dashboard icon to select an item</div>
          )}
        </div>
      );
    },
  }
);
