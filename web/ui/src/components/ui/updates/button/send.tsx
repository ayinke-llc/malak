import type { ServerAPIStatus } from "@/client/Api";
import { Button } from "@/components/Button";
import {
	Dialog,
	DialogClose,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "@/components/Dialog";
import { Label } from "@/components/Label";
import { Switch } from "@/components/Switch";
import client from "@/lib/client";
import { CREATE_CONTACT_MUTATION } from "@/lib/query-constants";
import { yupResolver } from "@hookform/resolvers/yup";
import { RiMailSendLine } from "@remixicon/react";
import { useMutation } from "@tanstack/react-query";
import type { AxiosError } from "axios";
import * as EmailValidator from "email-validator";
import { Option } from "lucide-react";
import { useState } from "react";
import { type SubmitHandler, useForm } from "react-hook-form";
import CreatableSelect from "react-select/creatable";
import { toast } from "sonner";
import * as yup from "yup";
import { Select } from "../../custom/select/select";
import type { ButtonProps } from "./props";

interface Option {
	readonly label: string;
	readonly value: string;
}

type SendUpdateInput = {
	email: string;
	link?: boolean;
	recipients?: Option[];
};

const schema = yup
	.object({
		email: yup.string().min(5).max(50).required(),
		link: yup.boolean().optional(),
	})
	.required();

const SendUpdateButton = ({}: ButtonProps) => {
	const [loading, setLoading] = useState<boolean>(false);

	const [options, setOptions] = useState<Option[]>([
		{ value: "oops", label: "oops" },
		{ value: "test", label: "test" },
	]);

	const [value, setValue] = useState<Option[]>([]);

	const contactMutation = useMutation({
		mutationKey: [CREATE_CONTACT_MUTATION],
		mutationFn: (data: { email: string }) =>
			client.contacts.contactsCreate(data),
		onSuccess: ({ data }) => {
			toast.info(`${data.contact.email} has been added as a contact now`);

			const newOption = {
				value: data.contact.id,
				label: data.contact.email,
			} as Option;

			setOptions((prev) => [...prev, newOption]);
			setValue((prev) => [...prev, newOption]);
		},
		onError(err: AxiosError<ServerAPIStatus>) {
			let msg = err.message;
			if (err.response !== undefined) {
				msg = err.response.data.message;
			}
			toast.error(msg);
		},
		retry: false,
		gcTime: Number.POSITIVE_INFINITY,
		onSettled: () => setLoading(false),
	});

	const createNewContact = (inputValue: string) => {
		setLoading(true);

		if (!EmailValidator.validate(inputValue)) {
			toast.error("you can only add an email address as a new recipient");
			setLoading(false);
			return;
		}

		contactMutation.mutate({ email: inputValue });
	};

	const {
		register,
		formState: { errors },
		handleSubmit,
	} = useForm({
		resolver: yupResolver(schema),
	});

	const onSubmit: SubmitHandler<SendUpdateInput> = (data) => {
		setLoading(true);
		console.log(data);
	};

	return (
		<>
			<div className="flex justify-center">
				<Dialog>
					<DialogTrigger asChild>
						<Button type="submit" size="lg" variant="primary" className="gap-1">
							<RiMailSendLine size={18} />
							Send
						</Button>
					</DialogTrigger>
					<DialogContent className="sm:max-w-lg">
						<form onSubmit={handleSubmit(onSubmit)}>
							<DialogHeader>
								<DialogTitle>Send this update</DialogTitle>
								<DialogDescription className="mt-1 text-sm leading-6">
									An email will be sent to all selected contacts immediately.
									Please re-verify your content is ready and good to go
								</DialogDescription>

								<div className="mt-4">
									<CreatableSelect
										isMulti
										isClearable
										isDisabled={loading}
										isLoading={loading}
										onChange={(newValue) => {
											setValue(newValue);
										}}
										onCreateOption={createNewContact}
										options={options}
										value={value}
									/>
								</div>

								<div className="mt-4">
									<Label htmlFor="select-author" className="font-medium">
										Select Author
									</Label>
									<Select
										data={[
											{
												label: "Lanre Adelowo",
												value: "lanre",
											},
											{
												label: "Ayinke",
												value: "ayinke",
											},
										]}
									/>
								</div>

								<div className="mt-4">
									<Switch disabled id="r3" {...register("link")} />
									<Label disabled htmlFor="r3">
										Coming soon. Generate a public viewable link for this update
									</Label>
								</div>
							</DialogHeader>
							<DialogFooter className="mt-6">
								<DialogClose asChild>
									<Button
										type={"button"}
										className="mt-2 w-full sm:mt-0 sm:w-fit"
										variant="secondary"
										isLoading={loading}
									>
										Cancel
									</Button>
								</DialogClose>
								<Button
									type="submit"
									className="w-full sm:w-fit"
									isLoading={loading}
								>
									Send
								</Button>
							</DialogFooter>
						</form>
					</DialogContent>
				</Dialog>
			</div>
		</>
	);
};

export default SendUpdateButton;
