set deleteTarget to "%s"
set foundMatch to false

tell application "System Events"
	tell application process "FaceTime"
		tell first window -- default window
			tell second scroll area -- unlabelled
				tell (first list where description is "Recent Calls")
					if ((first list where description is "Upcoming") exists) then
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
										repeat with checkText in (every static text as list)
											if (name of checkText is deleteTarget) then
												set foundMatch to true
												click button "Delete Link"
												exit repeat
											end if
										end repeat
									end tell
								end tell
							end repeat
						end tell
					end if
				end tell
			end tell
			if foundMatch then
				tell first sheet
					tell button "Delete Link"
						click
					end tell
				end tell
			end if
		end tell
	end tell
end tell
return foundMatch