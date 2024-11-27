export interface PopUpMessage {
	id?: number;
	pop_type?: number;
	title?: string;
	description?: string;
	enabled?: boolean;
	show_time?: number;
	close_time?: number;
	date_created?: string;
	popup_pos?: number;
	image_name?: string;
	isTrackable?: boolean;
}