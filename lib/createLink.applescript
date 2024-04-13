set copyLinkShareButtonSize to {216, 22}

do shell script "open facetime://"
tell application "System Events"
	tell application process "FaceTime"
		tell first window
			repeat until (button "Create Link" exists)
				delay 0.1
			end repeat
			tell button "Create Link"
				click
				repeat until (first pop over exists)
					delay 0.1
				end repeat
				tell first pop over
					tell first group
						repeat until (third button exists)
							delay 0.1
						end repeat
						set childElements to every UI element
						repeat with childElement in childElements
							if childElement's size is copyLinkShareButtonSize then
								click childElement
								delay 0.1
								return (the clipboard as text)
								exit repeat
							end if
						end repeat
					end tell
				end tell
			end tell
		end tell
	end tell
end tell