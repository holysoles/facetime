set copyLinkShareButtonSize to {216, 22}

tell application "System Events"
	tell application process "FaceTime"
		tell first window -- default window
			tell scroll area 2 -- unlabelled
				tell (first list where description is "Recent Calls")
					tell (first list where description is "Upcoming")
						tell first group -- most recently created call/link
							repeat with anAction in (actions as list)
								if (description of anAction) is "FaceTime Video" then
									set desiredAction to anAction
									exit repeat
								end if
							end repeat
							perform desiredAction
						end tell
					end tell
				end tell
			end tell
			tell button "Join"
				repeat until (exists)
					delay 0.1
				end repeat
				click
			end tell
		end tell
	end tell
end tell