// Extended the copy found at https://gist.github.com/ilkou/7bf2dbd42a7faf70053b43034fc4b5a4?permalink_comment_id=5194437#gistcomment-5194437
// This extended copy allows custom input amongst others
import React, { ReactElement, Ref, KeyboardEvent } from 'react';
import SelectComponent, {
  components,
  ClassNamesConfig,
  DropdownIndicatorProps,
  GroupBase,
  StylesConfig,
  MultiValueRemoveProps,
  ClearIndicatorProps,
  OptionProps,
  MenuProps,
  MenuListProps,
  Props,
  SelectInstance,
  createFilter,
} from 'react-select';
import CreatableSelect from 'react-select/creatable';
import { FixedSizeList as List } from 'react-window';
import { cn } from '@/lib/utils';
import { Check, ChevronDown, X } from 'lucide-react';

/** select option type */
export type OptionType = {
  label: string;
  value: string | number;
  data?: any;
  emails: {
    email: string
    reference: string
  }[];
};

// Extended props interface
export interface ExtendedSelectProps<IsMulti extends boolean = false>
  extends Omit<Props<OptionType, IsMulti>, 'onChange'> {
  options?: OptionType[];
  isMulti?: IsMulti;
  onCustomInputEnter?: (inputValue: string) => void;
  onChange?: (value: IsMulti extends true ? OptionType[] : OptionType | null) => void;
  allowCustomInput?: boolean;
}

/**
 * styles that aligns with shadcn/ui
 */
const selectStyles = {
  controlStyles: {
    base: 'flex !min-h-9 w-full rounded-md border border-input bg-transparent pl-3 py-1 pr-1 gap-1 text-sm shadow-sm transition-colors hover:cursor-pointer',
    focus: 'outline-none ring-1 ring-ring',
    disabled: 'cursor-not-allowed opacity-50',
  },
  placeholderStyles: 'text-muted-foreground text-sm ml-1 font-medium',
  valueContainerStyles: 'gap-1',
  multiValueStyles:
    'inline-flex items-center gap-2 rounded-md border border-transparent bg-secondary text-secondary-foreground hover:bg-secondary/80 px-1.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2',
  indicatorsContainerStyles: 'gap-1',
  clearIndicatorStyles: 'p-1 rounded-md',
  indicatorSeparatorStyles: 'bg-muted',
  dropdownIndicatorStyles: 'p-1 rounded-md',
  menu: 'mt-1.5 p-1.5 border border-input bg-background text-sm rounded-lg',
  menuList: 'morel-scrollbar',
  groupHeadingStyles:
    'py-2 px-1 text-secondary-foreground text-sm font-semibold',
  optionStyles: {
    base: 'hover:cursor-pointer hover:bg-accent hover:text-accent-foreground px-2 py-1.5 rounded-sm !text-sm !cursor-default !select-none !outline-none font-sans',
    focus: 'active:bg-accent/90 bg-accent text-accent-foreground',
    disabled: 'pointer-events-none opacity-50',
    selected: '',
  },
  noOptionsMessageStyles:
    'text-muted-foreground py-4 text-center text-sm border border-border rounded-sm',
  label: 'text-muted-foreground text-sm',
  loadingIndicatorStyles: 'flex items-center justify-center h-4 w-4 opacity-50',
  loadingMessageStyles: 'text-accent-foreground p-2 bg-accent',
};

/**
 * This factory method is used to build custom classNames configuration
 */
export const createClassNames = (
  classNames: ClassNamesConfig<OptionType, boolean, GroupBase<OptionType>>
): ClassNamesConfig<OptionType, boolean, GroupBase<OptionType>> => {
  return {
    clearIndicator: (state) =>
      cn(
        selectStyles.clearIndicatorStyles,
        classNames?.clearIndicator?.(state)
      ),
    container: (state) => cn(classNames?.container?.(state)),
    control: (state) =>
      cn(
        selectStyles.controlStyles.base,
        state.isDisabled && selectStyles.controlStyles.disabled,
        state.isFocused && selectStyles.controlStyles.focus,
        classNames?.control?.(state)
      ),
    dropdownIndicator: (state) =>
      cn(
        selectStyles.dropdownIndicatorStyles,
        classNames?.dropdownIndicator?.(state)
      ),
    group: (state) => cn(classNames?.group?.(state)),
    groupHeading: (state) =>
      cn(selectStyles.groupHeadingStyles, classNames?.groupHeading?.(state)),
    indicatorsContainer: (state) =>
      cn(
        selectStyles.indicatorsContainerStyles,
        classNames?.indicatorsContainer?.(state)
      ),
    indicatorSeparator: (state) =>
      cn(
        selectStyles.indicatorSeparatorStyles,
        classNames?.indicatorSeparator?.(state)
      ),
    input: (state) => cn(classNames?.input?.(state)),
    loadingIndicator: (state) =>
      cn(
        selectStyles.loadingIndicatorStyles,
        classNames?.loadingIndicator?.(state)
      ),
    loadingMessage: (state) =>
      cn(
        selectStyles.loadingMessageStyles,
        classNames?.loadingMessage?.(state)
      ),
    menu: (state) => cn(selectStyles.menu, classNames?.menu?.(state)),
    menuList: (state) => cn(classNames?.menuList?.(state)),
    menuPortal: (state) => cn(classNames?.menuPortal?.(state)),
    multiValue: (state) =>
      cn(selectStyles.multiValueStyles, classNames?.multiValue?.(state)),
    multiValueLabel: (state) => cn(classNames?.multiValueLabel?.(state)),
    multiValueRemove: (state) => cn(classNames?.multiValueRemove?.(state)),
    noOptionsMessage: (state) =>
      cn(
        selectStyles.noOptionsMessageStyles,
        classNames?.noOptionsMessage?.(state)
      ),
    option: (state) =>
      cn(
        selectStyles.optionStyles.base,
        state.isFocused && selectStyles.optionStyles.focus,
        state.isDisabled && selectStyles.optionStyles.disabled,
        state.isSelected && selectStyles.optionStyles.selected,
        classNames?.option?.(state)
      ),
    placeholder: (state) =>
      cn(selectStyles.placeholderStyles, classNames?.placeholder?.(state)),
    singleValue: (state) => cn(classNames?.singleValue?.(state)),
    valueContainer: (state) =>
      cn(
        selectStyles.valueContainerStyles,
        classNames?.valueContainer?.(state)
      ),
  };
};

export const defaultClassNames = createClassNames({});
export const defaultStyles: StylesConfig<
  OptionType,
  boolean,
  GroupBase<OptionType>
> = {
  input: (base) => ({
    ...base,
    'input:focus': {
      boxShadow: 'none',
    },
  }),
  multiValueLabel: (base) => ({
    ...base,
    whiteSpace: 'normal',
    overflow: 'visible',
  }),
  control: (base) => ({
    ...base,
    transition: 'none',
  }),
  menuList: (base) => ({
    ...base,
    '::-webkit-scrollbar': {
      background: 'transparent',
    },
    '::-webkit-scrollbar-track': {
      background: 'transparent',
    },
    '::-webkit-scrollbar-thumb': {
      background: 'hsl(var(--border))',
    },
    '::-webkit-scrollbar-thumb:hover': {
      background: 'transparent',
    },
  }),
};

/**
 * React select custom components
 */
export const DropdownIndicator = (
  props: DropdownIndicatorProps<OptionType>
) => {
  return (
    <components.DropdownIndicator {...props}>
      <ChevronDown className='h-4 w-4 opacity-50' />
    </components.DropdownIndicator>
  );
};

export const ClearIndicator = (props: ClearIndicatorProps<OptionType>) => {
  return (
    <components.ClearIndicator {...props}>
      <X className='h-4 w-4 opacity-50' />
    </components.ClearIndicator>
  );
};

export const MultiValueRemove = (props: MultiValueRemoveProps<OptionType>) => {
  return (
    <components.MultiValueRemove {...props}>
      <X className='h-3.5 w-3.5 opacity-50' />
    </components.MultiValueRemove>
  );
};

export const Option = (props: OptionProps<OptionType>) => {
  return (
    <components.Option {...props}>
      <div className='flex items-center justify-between'>
        <div>{props.label}</div>
        {props.isSelected && <Check className='h-4 w-4 opacity-50' />}
      </div>
    </components.Option>
  );
};

export const Menu = (props: MenuProps<OptionType>) => {
  return <components.Menu {...props}>{props.children}</components.Menu>;
};

export const MenuList = (props: MenuListProps<OptionType>) => {
  const { children, maxHeight } = props;
  const childrenArray = React.Children.toArray(children);

  const calculateHeight = () => {
    const totalHeight = childrenArray.length * 35;
    return totalHeight < maxHeight ? totalHeight : maxHeight;
  };

  const height = calculateHeight();

  if (!childrenArray || childrenArray.length - 1 === 0) {
    return <components.MenuList {...props} />;
  }

  return (
    <List
      height={height}
      itemCount={childrenArray.length}
      itemSize={35}
      width='100%'
    >
      {({ index, style }) => <div style={style}>{childrenArray[index]}</div>}
    </List>
  );
};

const Input = (props: any) => {
  const { onCustomInputEnter } = props.selectProps;

  const handleKeyDown = (event: KeyboardEvent<HTMLInputElement>) => {
    if (event.key === 'Enter' && props.value && onCustomInputEnter) {
      event.preventDefault();
      event.stopPropagation();
      onCustomInputEnter(props.value);
      // Reset only the input value using the selectProps inputValue
      props.selectProps.onInputChange('', { action: 'input-change' });
    }
  };

  return (
    <components.Input
      {...props}
      onKeyDown={(e: KeyboardEvent<HTMLInputElement>) => {
        handleKeyDown(e);
        props.onKeyDown?.(e);
      }}
    />
  );
};

const BaseSelect = <IsMulti extends boolean = false>(
  props: ExtendedSelectProps<IsMulti>,
  ref: React.Ref<SelectInstance<OptionType, IsMulti, GroupBase<OptionType>>>
) => {
  const {
    styles = defaultStyles,
    classNames = defaultClassNames,
    components = {},
    options = [],
    allowCustomInput = false,
    onCustomInputEnter,
    ...rest
  } = props;

  const instanceId = React.useId();

  const SelectVariant = allowCustomInput ? CreatableSelect : SelectComponent;

  return (
    <SelectVariant<OptionType, IsMulti, GroupBase<OptionType>>
      ref={ref}
      instanceId={instanceId}
      unstyled
      filterOption={createFilter({
        matchFrom: 'any',
        stringify: (option) => option.label,
      })}
      components={{
        DropdownIndicator,
        ClearIndicator,
        MultiValueRemove,
        Option,
        Menu,
        MenuList,
        Input,
        ...components,
      }}
      styles={styles}
      classNames={classNames}
      options={options}
      onCustomInputEnter={onCustomInputEnter}
      {...rest}
    />
  );
};

export default React.forwardRef(BaseSelect) as <
  IsMulti extends boolean = false,
>(
  p: ExtendedSelectProps<IsMulti> & {
    ref?: Ref<
      React.LegacyRef<
        SelectInstance<OptionType, IsMulti, GroupBase<OptionType>>
      >
    >;
  }
) => ReactElement;
