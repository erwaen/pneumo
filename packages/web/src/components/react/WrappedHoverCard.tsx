import { CalendarIcon } from "lucide-react"

import {
    Avatar,
    AvatarFallback,
    AvatarImage,
} from "@/components/ui/avatar"
import { Button } from "@/components/ui/button"
import {
    HoverCard,
    HoverCardContent,
    HoverCardTrigger,
} from "@/components/ui/hover-card"


export default function WrappedHoverCard({ title }: { title: string }) {
    return (
        <HoverCard>
            <HoverCardTrigger asChild>
                <Button className="m-0 p-0" variant="link">{title}</Button>
            </HoverCardTrigger>
            <HoverCardContent className="w-80">
                <div className="flex justify-between gap-1">
                    <Avatar>
                        <AvatarImage src="https://erwaen.github.io" />
                        <AvatarFallback className="bg-blue-100">EW</AvatarFallback>
                    </Avatar>
                    <div >
                        <h4 className="text-sm font-semibold">@erwaen</h4>
                        <p className="text-sm">
                            Erik Wasmosy – Portfolio
                        </p>
                    </div>
                </div>
            </HoverCardContent>
        </HoverCard>
    )
}
