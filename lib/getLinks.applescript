set copyLinkShareButtonSize to {216, 22}
set linkList to {}

do shell script "open facetime://"
tell application "System Events"
	tell application process "FaceTime"
		tell first window -- default window
			tell second scroll area -- unlabelled
				tell (first list where description is "Recent Calls")
					if (first list where description is "Upcoming" exists)
						tell (first list where description is "Upcoming")
							repeat with linkEntry in every group
								tell linkEntry
									repeat with anAction in (actions as list)
										if (description of anAction) is "Info" then
											set desiredAction to anAction
											exit repeat
										end if
									end repeat
									perform desiredAction
									repeat until (first pop over exists)
										delay 0.1
									end repeat
									tell first pop over
										set linkList to (linkList & (name of (first static text where name contains "facetime.apple.com")))
									end tell
								end tell
							end repeat
						end tell
					end if
				end tell
			end tell
		end tell
	end tell
end tell
return linkList