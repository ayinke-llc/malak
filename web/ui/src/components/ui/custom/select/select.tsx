import {
	Select as BaseSelect,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@/components/Select";
import type { RemixiconComponentType } from "@remixicon/react";
import { Avatar, AvatarFallback, AvatarImage } from "../avatar/avatar";

export type SelectProps = {
	data: {
		value: string;
		label: string;
		icon?: RemixiconComponentType;
		avatarURL?: string;
	}[];
};

export function Select({ data }: SelectProps) {
	return (
		<BaseSelect>
			<SelectTrigger>
				<SelectValue placeholder="Select" />
			</SelectTrigger>
			<SelectContent>
				{data.map((item) => (
					<SelectItem key={item.value} value={item.value}>
						<div className="flex items-center">
							{item.icon && <item.icon className="mr-2 h-4 w-4" />}
							{item.avatarURL && (
								<Avatar className="mr-2 h-4 w-4">
									<AvatarImage src={item.avatarURL} />
									<AvatarFallback>AY</AvatarFallback>
								</Avatar>
							)}
							{item.label}
						</div>
					</SelectItem>
				))}
			</SelectContent>
		</BaseSelect>
	);
}
