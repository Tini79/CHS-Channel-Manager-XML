procedure TCHSChannelManager.ProcessOTA_ResRetrieveRS(
  XMLString: string);
var
  Count1, Count2, Count3, Count4, Count5, Count6, Count7, Count8, CountAttribute,
  CountHotelReservation, CountHotelReservationField, CountRoomStays, CountRoomStay, CountCustomer,
  RoomStayID, Adult, Child: Integer;
  IsXMLSuccess, IsNoRPH: Boolean;
  POSX, BookingCode, OTAID, RPH, RoomNumber, RoomTypeCode, BedTypeCode, RoomRateCode, RoomRateAmountStr, AdultStr, ChildStr, InfantStr, ArrivalDateStr, DepartureDateStr, ResGuestRPH, ArrivalTimeStr, FullName, GivenName, MiddleName, Surname, Phone1, Email, Street, City, PostalCode, State, Country, Company,
  GuestProfileID, ContactPersonID, GuestDetailID, ResStatus: string;
  ParameterCondition: Int64;
  ArrivalDate, DepartureDate, ServerDate: TDateTime;
  ArrivalTime: TTime;
  RoomRateAmount: Double;
  RoomList: TStringList;
begin
  try
    RoomList := TStringList.Create;
    ProgramConfiguration.CCMSReservationAsAllotment := ReadConfigurationBoolean(SystemCode.Hotel, ConfigurationCategory.ServiceCCMS, ConfigurationName.CCMSSMReservationAsAllotment, False);
    if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
    begin
      DebugLog.Text := XMLString;
      DebugLog.SaveToFile('C:\Temp\XMLLog' +IntToStr(ProgramVariable.DebugCount)+ '.txt');
      Inc(ProgramVariable.DebugCount);
    end;

    DebugLog.Clear;
    DebugLogDetail.Clear;
    IsXMLSuccess := False;
    XMLDocument2.LoadFromXML(XMLString);
    for Count1 := 0 to XMLDocument2.DocumentElement.ChildNodes.Count - 1 do
    begin
      if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
        DebugLogDetail.Add(XMLDocument2.DocumentElement.ChildNodes[Count1].NodeName);
      if XMLDocument2.DocumentElement.ChildNodes[Count1].NodeName = 'Success' then
        //Success
        IsXMLSuccess := True
      else if IsXMLSuccess and (XMLDocument2.DocumentElement.ChildNodes[Count1].NodeName = 'ReservationsList') then
      begin
        //Reservation List
        for CountHotelReservation := 0 to XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes.Count - 1 do
        begin
          if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
            DebugLogDetail.Add(XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].NodeName);

          POSX := '';
          BookingCode := '';
          OTAID := '';
          XMLRoomStayList.Clear;
          XMLRPHList.Clear;
          XMLResGuestList.Clear;
          if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].NodeName = 'HotelReservation' then
          begin
            //Hotel Reservation
            if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
            begin
              DebugLog.Add('----------------------------------------------------------------------------------------------------');
              DebugLog.Add('Reservation: ' + IntToStr(CountHotelReservation + 1));
              DebugLog.Add('----------------------------------------------------------------------------------------------------');
            end;

            for CountAttribute := 0 to XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].AttributeNodes.Count - 1 do
            begin
              if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].AttributeNodes[CountAttribute].NodeName = 'ResStatus' then
              begin
                ResStatus := XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].AttributeNodes[CountAttribute].Text;
                if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                  DebugLog.Add('ResStatus :' + ResStatus);
              end;
            end;

            for CountHotelReservationField := 0 to XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes.Count - 1 do
            begin
              if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                DebugLogDetail.Add(XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].NodeName);

              if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].NodeName = 'POS' then
              begin
                //POS
                for Count2 := 0 to XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes.Count - 1 do
                begin
                  if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                    DebugLogDetail.Add(XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].NodeName);

                  if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].NodeName = 'Source' then
                  begin
                    //Source
                    for Count3 := 0 to XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes.Count - 1 do
                    begin
                      if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                        DebugLogDetail.Add(XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].NodeName);

                      if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].NodeName = 'BookingChannel' then
                      begin
                        //Booking Channel
                        for Count4 := 0 to XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes.Count - 1 do
                        begin
                          if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                            DebugLogDetail.Add(XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].NodeName);
                          if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].NodeName = 'CompanyName' then
                          begin
                            //CompanyName OTA
                            for CountAttribute := 0 to XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].AttributeNodes.Count - 1 do
                            begin
                              if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].AttributeNodes[CountAttribute].NodeName = 'Code' then
                              begin
                                //OTA Code
                                POSX := XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].AttributeNodes[CountAttribute].Text;
                                if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                                  DebugLog.Add('POS :' + POSX);

                                Break;
                              end;
                            end;
                            Break;
                          end;
                        end;
                        Break;
                      end;
                    end;
                    Break;
                  end;
                end;
              end
              else if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].NodeName = 'UniqueID' then
              begin
                //Unique ID
                if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].Attributes['Type'] = '14' then
                begin
                  //Unique ID 14 = Site Minder ID
                  BookingCode := XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].Attributes['ID'];
                  if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                    DebugLog.Add('Site Minder ID: ' + XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].Attributes['ID']);
                end
                else if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].Attributes['Type'] = '16' then
                begin
                  //Unique ID 16 = OTA ID
                  OTAID := XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].Attributes['ID'];
                  if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                    DebugLog.Add('OTA ID: ' + OTAID);
                end;
              end
              else if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].NodeName = 'RoomStays' then
              begin
                //Room Stays
                for CountRoomStays := 0 to XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes.Count - 1 do
                begin
                  if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                    DebugLogDetail.Add(XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].NodeName);

                  if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].NodeName = 'RoomStay' then
                  begin
                    //Room Stay
                    if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                    begin
                      DebugLog.Add('-------------------');
                      DebugLog.Add('Room Stay: ' + IntToStr(CountRoomStays + 1));
                      DebugLog.Add('-------------------');
                    end;

                    IsNoRPH := True;

                    RoomTypeCode := XMLEmptyData;
                    RoomRateCode := XMLEmptyData;
                    RoomRateAmountStr := '0';
                    AdultStr := XMLEmptyData;
                    ChildStr := XMLEmptyData;
                    InfantStr := XMLEmptyData;
                    ArrivalDateStr := XMLEmptyData;
                    DepartureDateStr := XMLEmptyData;

                    for CountRoomStay := 0 to XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes.Count - 1 do
                    begin
                      if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                        DebugLogDetail.Add(XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].NodeName);
                      if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].NodeName = 'RoomRates' then
                      begin
                        //Room Rates
                        for Count2 := 0 to XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes.Count - 1 do
                        begin
                          if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                            DebugLogDetail.Add(XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].NodeName);
                          if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].NodeName = 'RoomRate' then
                          begin
                            //Room Rate
                            for CountAttribute := 0 to XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].AttributeNodes.Count - 1 do
                            begin
                              if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                                DebugLogDetail.Add(XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].AttributeNodes[CountAttribute].NodeName);
                              if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].AttributeNodes[CountAttribute].NodeName = 'RatePlanCode' then
                              begin
                                //Room Rate Code
                                RoomRateCode := XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].AttributeNodes[CountAttribute].Text;
                                if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                                  DebugLog.Add('RoomRateCode: ' + RoomRateCode);
                              end
                              else if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].AttributeNodes[CountAttribute].NodeName = 'RoomTypeCode' then
                              begin
                                //Room Type Code
                                RoomTypeCode := XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].AttributeNodes[CountAttribute].Text;
                                if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                                  DebugLog.Add('RoomTypeCode: ' + RoomTypeCode);
                              end;
                            end;

                            for Count3 := 0 to XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].ChildNodes.Count - 1 do
                            begin
                              if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                                DebugLogDetail.Add(XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].ChildNodes[Count3].NodeName);
                              if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].ChildNodes[Count3].NodeName = 'Rates' then
                              begin
                                //Rates
                                for Count4 := 0 to XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].ChildNodes[Count3].ChildNodes.Count - 1 do
                                begin
                                  if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                                    DebugLogDetail.Add(XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].NodeName);
                                  if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].NodeName = 'Rate' then
                                  begin
                                    //Rate
                                    for Count5 := 0 to XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes.Count - 1 do
                                    begin
                                      if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                                        DebugLogDetail.Add(XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].NodeName);
                                      if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].NodeName = 'Total' then
                                      begin
                                        //Total
                                        for CountAttribute := 0 to XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].AttributeNodes.Count - 1 do
                                        begin
                                          if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                                            DebugLogDetail.Add(XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].AttributeNodes[CountAttribute].NodeName);
                                          if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].AttributeNodes[CountAttribute].NodeName = 'AmountAfterTax' then
                                          begin
                                            //Rate Amount After Tax
                                            RoomRateAmountStr := XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].AttributeNodes[CountAttribute].Text;
                                            if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                                              DebugLog.Add('RoomRateAmount: ' + RoomRateAmountStr);
                                            Break;
                                          end;
                                        end;
                                        Break;
                                      end;
                                    end;
                                    Break;
                                  end;
                                end;
                              end;
                            end;
                          end;
                        end;
                      end
                      else if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].NodeName = 'GuestCounts' then
                      begin
                        //GuestCounts
                        for Count2 := 0 to XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes.Count - 1 do
                        begin
                          if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                            DebugLogDetail.Add(XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].NodeName);
                          if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].NodeName = 'GuestCount' then
                          begin
                            //Guest Count
                            if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                              DebugLogDetail.Add('Count Attr Guest Count ' + IntToStr(Count2 + 1) + ': ' + IntToStr(XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].AttributeNodes.Count));
                            if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].AttributeNodes.Count > 1 then
                            begin
                              //Adult
                              if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                                DebugLogDetail.Add('Age Qualifying: ' + XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].Attributes['AgeQualifyingCode']);
                              if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].Attributes['AgeQualifyingCode'] = '10' then
                              begin
                                AdultStr := XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].Attributes['Count'];
                                if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                                  DebugLog.Add('Adult: ' + AdultStr);
                              end
                              //Child
                              else if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].Attributes['AgeQualifyingCode'] = '8' then
                              begin
                                ChildStr := XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].Attributes['Count'];
                                if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                                  DebugLog.Add('Child: ' + ChildStr);
                              end
                              //Infant
                              else if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].Attributes['AgeQualifyingCode'] = '7' then
                              begin
                                InfantStr := XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].Attributes['Count'];
                                if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                                  DebugLog.Add('Infant: ' + InfantStr);
                              end;
                            end;
                          end;
                        end;
                      end
                      else if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].NodeName = 'TimeSpan' then
                      begin
                        //TimeSpan
                        for CountAttribute := 0 to XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].AttributeNodes.Count - 1 do
                        begin
                          if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].AttributeNodes[CountAttribute].NodeName = 'Start' then
                          begin
                            //Arrival Date
                            ArrivalDateStr := XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].AttributeNodes[CountAttribute].Text;
                            if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                              DebugLog.Add('Arrival Date: ' + ArrivalDateStr);
                          end
                          else if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].AttributeNodes[CountAttribute].NodeName = 'End' then
                          begin
                            //Depature Date
                            DepartureDateStr := XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].AttributeNodes[CountAttribute].Text;
                            if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                              DebugLog.Add('Departure Date: ' + DepartureDateStr);
                          end;
                        end;
                        //------------------------------------------------------------------------------------------------------------------------
                        //Save Room Stay
                        //------------------------------------------------------------------------------------------------------------------------
                        if (RoomTypeCode <> '') and (RoomRateCode <> '') and (ArrivalDateStr <> '') and (DepartureDateStr <> '') then
                        beginti
                          try
                            if AdultStr = '' then
                              AdultStr := '1';
                            if ChildStr = '' then
                              AdultStr := '0';
                            if InfantStr = '' then
                              AdultStr := '0';

                            XMLRoomStayList.Add(RoomTypeCode + DelimiterX + RoomRateCode + DelimiterX + AdultStr + DelimiterX + ChildStr + DelimiterX + InfantStr + DelimiterX + ArrivalDateStr + DelimiterX + DepartureDateStr);
                          except
                          end;
                        end;
                      end
                      else if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].NodeName = 'ResGuestRPHs' then
                      begin
                        //ResGuestRPHs
                        for Count2 := 0 to XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes.Count - 1 do
                        begin
                          if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                            DebugLogDetail.Add(XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].NodeName);
                          if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].NodeName = 'ResGuestRPH' then
                          begin
                            if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                              DebugLogDetail.Add('Count Attr Time Span ' + IntToStr(XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].AttributeNodes.Count));
                            for CountAttribute := 0 to XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].AttributeNodes.Count - 1 do
                            begin
                              if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].AttributeNodes[CountAttribute].NodeName = 'RPH' then
                              begin
                                //RPH
                                if XMLRoomStayList.Count > 0 then
                                begin
                                  IsNoRPH := False;
                                  XMLRPHList.Add(IntToStr(XMLRoomStayList.Count - 1) + DelimiterX + XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[CountRoomStays].ChildNodes[CountRoomStay].ChildNodes[Count2].AttributeNodes[CountAttribute].Text);
                                  if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                                    DebugLog.Add('RPH ' + IntToStr(XMLRPHList.Count) + ': ' + XMLRPHList.Strings[XMLRPHList.Count - 1]);
                                end;
                              end;
                            end;
                          end;
                        end;
                      end;
                    end;
                    //------------------------------------------------------------------------------------------------------------------------
                    //Check if No RPH
                    //------------------------------------------------------------------------------------------------------------------------
                    if IsNoRPH then
                    begin
                      if XMLRoomStayList.Count > 0 then
                      begin
                        XMLRPHList.Add(IntToStr(XMLRoomStayList.Count - 1) + DelimiterX + '1');
                        if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                          DebugLog.Add('Manual RPH ' + IntToStr(XMLRPHList.Count) + ': ' + XMLRPHList.Strings[XMLRPHList.Count - 1]);
                      end;
                    end;
                  end;
                end;
              end
              else if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].NodeName = 'ResGuests' then
              begin
                //ResGuests
                for Count2 := 0 to XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes.Count - 1 do
                begin
                  if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                    DebugLogDetail.Add(XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].NodeName);

                  if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].NodeName = 'ResGuest' then
                  begin
                    //ResGuests
                    if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                    begin
                      DebugLog.Add('---------------------------------------------------------------');
                      DebugLog.Add('Res Guest ' + IntToStr(Count2 +  1));
                      DebugLog.Add('---------------------------------------------------------------');
                      DebugLogDetail.Add('Count Res Guest Attr: ' + IntToStr(XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].AttributeNodes.Count));
                    end;

                    ResGuestRPH := XMLEmptyData;
                    ArrivalTimeStr := XMLEmptyData;
                    GivenName := XMLEmptyData;
                    MiddleName := XMLEmptyData;
                    Surname := XMLEmptyData;
                    Phone1 := XMLEmptyData;
                    Email := XMLEmptyData;
                    Street := XMLEmptyData;
                    City := XMLEmptyData;
                    PostalCode := XMLEmptyData;
                    State := XMLEmptyData;
                    Country := XMLEmptyData;
                    Company := XMLEmptyData;

                    for CountAttribute := 0 to XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].AttributeNodes.Count - 1 do
                    begin
                      if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                        DebugLogDetail.Add(XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].AttributeNodes[CountAttribute].NodeName);
                      if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].AttributeNodes[CountAttribute].NodeName = 'ResGuestRPH' then
                      begin
                        //ResGuestRPH
                        ResGuestRPH := XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].AttributeNodes[CountAttribute].Text;
                        if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                          DebugLog.Add('ResGuestRPH: ' + ResGuestRPH);
                      end
                      else if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].AttributeNodes[CountAttribute].NodeName = 'ArrivalTime' then
                      begin
                        //Arrival Time
                        ArrivalTimeStr := XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].AttributeNodes[CountAttribute].Text;
                        if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                          DebugLog.Add('Arrival Time: ' + ArrivalTimeStr);
                      end;
                    end;

                    for Count3 := 0 to XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes.Count - 1 do
                    begin
                      if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                        DebugLogDetail.Add(XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].NodeName);
                      if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].NodeName = 'Profiles' then
                      begin
                        //Profiles
                        for Count4 := 0 to XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes.Count - 1 do
                        begin
                          if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                            DebugLogDetail.Add(XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].NodeName);
                          if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].NodeName = 'ProfileInfo' then
                          begin
                            //Profile Info
                            for Count5 := 0 to XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes.Count - 1 do
                            begin
                              if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                                DebugLogDetail.Add(XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].NodeName);
                              if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].NodeName = 'Profile' then
                              begin
                                //Profile
                                for Count6 := 0 to XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes.Count - 1 do
                                begin
                                  if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                                    DebugLogDetail.Add(XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].NodeName);
                                  if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].NodeName = 'Customer' then
                                  begin
                                    //Customer
                                    for CountCustomer := 0 to XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].ChildNodes.Count - 1 do
                                    begin
                                      if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                                        DebugLogDetail.Add(XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].ChildNodes[CountCustomer].NodeName);
                                      if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].ChildNodes[CountCustomer].NodeName = 'PersonName' then
                                      begin
                                        for Count7 := 0 to XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].ChildNodes[CountCustomer].ChildNodes.Count - 1 do
                                        begin
                                          if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].ChildNodes[CountCustomer].ChildNodes[Count7].NodeName = 'GivenName' then
                                          begin
                                            //Given Name
                                            GivenName := XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].ChildNodes[CountCustomer].ChildNodes[Count7].Text;
                                            if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                                              DebugLog.Add('Given Name: ' + GivenName);
                                          end
                                          else if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].ChildNodes[CountCustomer].ChildNodes[Count7].NodeName = 'MiddleName' then
                                          begin
                                            //Middle Name
                                            MiddleName := XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].ChildNodes[CountCustomer].ChildNodes[Count7].Text;
                                            if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                                              DebugLog.Add('Middle Name: ' + MiddleName);
                                          end
                                          else if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].ChildNodes[CountCustomer].ChildNodes[Count7].NodeName = 'Surname' then
                                          begin
                                            //Surname
                                            Surname := XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].ChildNodes[CountCustomer].ChildNodes[Count7].Text;
                                            if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                                              DebugLog.Add('Surname: ' + Surname);
                                          end
                                        end;
                                      end
                                      else if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].ChildNodes[CountCustomer].NodeName = 'Telephone' then
                                      begin
                                        for CountAttribute := 0 to XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].ChildNodes[CountCustomer].AttributeNodes.Count - 1 do
                                        begin
                                          if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].ChildNodes[CountCustomer].AttributeNodes[CountAttribute].NodeName = 'PhoneNumber' then
                                          begin
                                            //Phone Number
                                            Phone1 := XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].ChildNodes[CountCustomer].AttributeNodes[CountAttribute].Text;
                                            if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                                              DebugLog.Add('Phone Number: ' + Phone1);
                                          end;
                                        end;
                                      end
                                      else if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].ChildNodes[CountCustomer].NodeName = 'Email' then
                                      begin
                                        //Email
                                        Email := XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].ChildNodes[CountCustomer].Text;
                                        if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                                          DebugLog.Add('Email: ' + Email);
                                      end
                                      else if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].ChildNodes[CountCustomer].NodeName = 'Address' then
                                      begin
                                        for Count7 := 0 to XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].ChildNodes[CountCustomer].ChildNodes.Count - 1 do
                                        begin
                                          if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].ChildNodes[CountCustomer].ChildNodes[Count7].NodeName = 'AddressLine' then
                                          begin
                                            //Address
                                            if Street = XMLEmptyData then
                                              Street := XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].ChildNodes[CountCustomer].ChildNodes[Count7].Text
                                            else
                                              Street := Street + ' ' + XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].ChildNodes[CountCustomer].ChildNodes[Count7].Text;
                                            if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                                              DebugLog.Add('Street: ' + Street);
                                          end
                                          else if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].ChildNodes[CountCustomer].ChildNodes[Count7].NodeName = 'CityName' then
                                          begin
                                            //City Name
                                            City := XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].ChildNodes[CountCustomer].ChildNodes[Count7].Text;
                                            if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                                              DebugLog.Add('City Name: ' + City);
                                          end
                                          else if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].ChildNodes[CountCustomer].ChildNodes[Count7].NodeName = 'PostalCode' then
                                          begin
                                            //Postal Code
                                            PostalCode := XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].ChildNodes[CountCustomer].ChildNodes[Count7].Text;
                                            if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                                              DebugLog.Add('Postal Code: ' + PostalCode);
                                          end
                                          else if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].ChildNodes[CountCustomer].ChildNodes[Count7].NodeName = 'StateProv' then
                                          begin
                                            //State Prov
                                            State := XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].ChildNodes[CountCustomer].ChildNodes[Count7].Text;
                                            if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                                              DebugLog.Add('State: ' + State);
                                          end
                                          else if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].ChildNodes[CountCustomer].ChildNodes[Count7].NodeName = 'CountryName' then
                                          begin
                                            //CountryName
                                            Country := XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].ChildNodes[CountCustomer].ChildNodes[Count7].Text;
                                            if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                                              DebugLog.Add('Country: ' + Country);
                                          end
                                          else if XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].ChildNodes[CountCustomer].ChildNodes[Count7].NodeName = 'CompanyName' then
                                          begin
                                            //CompanyName
                                            Company := XMLDocument2.DocumentElement.ChildNodes[Count1].ChildNodes[CountHotelReservation].ChildNodes[CountHotelReservationField].ChildNodes[Count2].ChildNodes[Count3].ChildNodes[Count4].ChildNodes[Count5].ChildNodes[Count6].ChildNodes[CountCustomer].ChildNodes[Count7].Text;
                                            if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                                              DebugLog.Add('Company: ' + Company);
                                          end;
                                        end;
                                      end;
                                    end;
                                  end
                                end;
                              end;
                            end;
                          end;
                        end;
                      end;
                    end;
                    //------------------------------------------------------------------------------------------------------------------------
                    //Save Res Guest
                    //------------------------------------------------------------------------------------------------------------------------
                    if (ResGuestRPH <> '') and (GivenName <> '') then
                      XMLResGuestList.Add(ResGuestRPH + DelimiterX + ArrivalTimeStr + DelimiterX + GivenName + DelimiterX + MiddleName + DelimiterX + Surname + DelimiterX +  Phone1 + DelimiterX + Email + DelimiterX + Street + DelimiterX + City + DelimiterX + PostalCode + DelimiterX + State + DelimiterX + Country + DelimiterX + Company);
                  end;
                end;
              end;
            end;
          end;

          //------------------------------------------------------------------------------------------------------------------------
          //Save Log
          //------------------------------------------------------------------------------------------------------------------------
          if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
          begin
            XMLRoomStayList.SaveToFile('C:\Temp\XMLRoomStayList' +IntToStr(CountHotelReservation)+ '.txt');
            XMLRPHList.SaveToFile('C:\Temp\XMLRPHList' +IntToStr(CountHotelReservation)+ '.txt');
            XMLResGuestList.SaveToFile('C:\Temp\XMLResGuestList' +IntToStr(CountHotelReservation)+ '.txt');
          end;
          //------------------------------------------------------------------------------------------------------------------------
          //Insert to Reservation
          //------------------------------------------------------------------------------------------------------------------------
          if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
            DebugLogDetail.Add('Booking Code: ' + BookingCode + 'OTA ID: ' + OTAID);
          if ResStatus = 'Book' then
          begin
            if not IsCMReservationReceived(BookingCode, OTAID) then
            begin
              for Count2 := 0 to XMLRPHList.Count - 1 do == (digunakan atau tidak)
              begin
                DeployRPH(XMLRPHList.Strings[Count2], RoomStayID, RPH);
                if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                  DebugLogDetail.Add('InsertRV, RPH: ' + RPH);
                if RPH <> XMLEmptyData then
                begin
                  DeployRoomStay(XMLRoomStayList.Strings[RoomStayID], RoomTypeCode, BedTypeCode, RoomRateCode, Adult, Child, ArrivalDate, DepartureDate);
                  // Contohnya nanti dari booknlink itu seperti ini STD#DOBL ini artinya roomtype STD bed type DOBL
                  if BedTypeCode = '' then
                  begin
                    RoomList.Text := GetAvailableRoomByType(RoomTypeCode, '', FormatDateTimeX(ArrivalDate), FormatDateTimeX(DepartureDate), 0, 0, 0, 0, False, ProgramConfiguration.CCMSReservationAsAllotment);
                    if RoomList.Count > 0 then
                      BedTypeCode := GetBedTypeCode(RoomList.Strings[0]);
                  end;

                  // if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                  //   DebugLogDetail.Add('DeployRoomStay: ' + ':' + RoomTypeCode + ':' + BedTypeCode + ':' + RoomRateCode + ':' + RoomRateAmountStr + ':' + IntToStr(Adult) + ':' + IntToStr(Child) + ':' + FormatDateTimeX2(ArrivalDate) + ':' + FormatDateTimeX2(DepartureDate));

                  for Count3 := 0 to XMLResGuestList.Count - 1 do
                  // Tanyakan
                  begin
                    DeployResGuest(XMLResGuestList.Strings[Count3], ArrivalTime, ResGuestRPH, FullName, Phone1, Email, Street, City, PostalCode, State, Country, Company);
                    FullName := UpperCase(FullName);
                    FullName := TRIM(ReplaceText(UpperCase(FullName), 'NO SURNAME', ''));

                    // if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                    //   DebugLogDetail.Add('RPH = ResRPH ' + RPH + ':' + ResGuestRPH);
                    if RPH = ResGuestRPH then
                    // tanyakan RPH
                    begin
                      ServerDate := GetServerDate;
                      ProgramVariable.AuditDate := GetAuditDate;
                      ProgramConfiguration.CheckOutLimit := StrToTime(ReadConfigurationString(SystemCode.Hotel, ConfigurationCategory.Reservation, ConfigurationName.CheckOutLimit, False), ProgramVariable.FormatSettingX);
                      ReplaceTime(ArrivalDate, ArrivalTime);
                      ReplaceTime(DepartureDate, ProgramConfiguration.CheckOutLimit);

                      if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                      begin
                        DebugLogDetail.Add('Check 1: ' + BoolToStr((DateOf(ArrivalDate) >= DateOf(ProgramVariable.AuditDate))));
                        DebugLogDetail.Add('Check 2: ' + BoolToStr((DateOf(DepartureDate) > DateOf(ProgramVariable.AuditDate))));
                        DebugLogDetail.Add('Check 3: ' + BoolToStr((GetAvailableRoomCountByType(RoomTypeCode, BedTypeCode, ArrivalDate, DepartureDate, 0, 0, 0, 0, False, True) > 0)));
                        DebugLogDetail.Add('OTA ID: ' + OTAID);
                      end;

                      RoomRateAmount := 0;
                      if (RoomRateAmountStr <> '') and (RoomRateAmountStr <> '0') then
                      begin
                        try
                          RoomRateAmount := StrToFloat(RoomRateAmountStr, FormatStringTemp);
                        except
                        end;
                      end;

                      if (DateOf(ArrivalDate) >= DateOf(ProgramVariable.AuditDate)) and (DateOf(DepartureDate) > DateOf(ProgramVariable.AuditDate)) and
                         (GetAvailableRoomCountByType(RoomTypeCode, BedTypeCode, ArrivalDate, DepartureDate, 0, 0, 0, 0, False, ProgramConfiguration.CCMSReservationAsAllotment) > 0) then
                      begin

                        GuestProfileID := InsertGuestProfile('', FullName, Street, '', City, '', '', '', PostalCode, Phone1, '', '', Email, '', '', '', '', '', BoolToStr(True), '',  '0000-00-00', CPType.Guest,
                                                             '', '', '', '', '', '', '', '', '', '', '', '',
                                                             '', '', '', '', '', '', '', '', '', '', '', '',
                                                             '', GuestProfileSource.Hotel, ServerDate);

                        ContactPersonID := InsertContactPerson('', FullName, Street, '', City, '', '', '', PostalCode, Phone1, '', '', Email, '', '', '', '', '', BoolToStr(True), '',  '0000-00-00', CPType.Guest,
                                                             '', '', '', '', '', '', '', '', '', '', '', '',
                                                             '', '', '', '', '', '', '', '', '', '', '', '');

                        GuestDetailID := InsertGuestDetail(ArrivalDate, DepartureDate, Adult, Child, RoomTypeCode, BedTypeCode, '', RoomRateCode, '', '', '', '', '', '', True, RoomRateAmount, RoomRateAmount, 0, 0);
                        ParameterCondition :=   (ContactPersonID, '', '', '', GuestDetailID, '', GuestProfileID, '', '', '', '', FullName, '', '', ReservationStatus.New, '', '', '', '', '', '', BookingCode, OTAID, ResStatus, 1, NullDate, NullDate, True, False, ProgramConfiguration.CCMSReservationAsAllotment);
                        AssignRoom(ParameterCondition, False);
                       end
                      else
                      begin
                        GuestProfileID := InsertGuestProfile('', FullName, Street, '', City, '', '', '', PostalCode, Phone1, '', '', Email, '', '', '', '', '', BoolToStr(True), '',  '0000-00-00', CPType.Guest,
                                                             '', '', '', '', '', '', '', '', '', '', '', '',
                                                             '', '', '', '', '', '', '', '', '', '', '', '',
                                                             '', GuestProfileSource.Hotel, ServerDate);

                        ContactPersonID := InsertContactPerson('', FullName, Street, '', City, '', '', '', PostalCode, Phone1, '', '', Email, '', '', '', '', '', BoolToStr(True), '',  '0000-00-00', CPType.Guest,
                                                             '', '', '', '', '', '', '', '', '', '', '', '',
                                                             '', '', '', '', '', '', '', '', '', '', '', '');

                        GuestDetailID := InsertGuestDetail(ArrivalDate, DepartureDate, Adult, Child, RoomTypeCode, BedTypeCode, '', RoomRateCode, '', '', '', '', '', '', True, RoomRateAmount, RoomRateAmount, 0, 0);
                        ParameterCondition := InsertReservation(ContactPersonID, '', '', '', GuestDetailID, '', GuestProfileID, '', '', '', '', FullName, '', '', ReservationStatus.WaitList, '', '', '', '', '', '', BookingCode, OTAID, ResStatus, 1, NullDate, NullDate, True, False, ProgramConfiguration.CCMSReservationAsAllotment);
                      end;
                      //                    InsertLogUser(LogUserAction.InsertReservation, IntToStr(ParameterCondition), '', '', '', LogUserAction.InsertReservationX);
                      //    //          ProcessSMSSchedule(SMSevent.OnInsertReservation, 'reservation.number = "' +IntToStr(ParameterCondition)+ '"', '', '', '', '', '', '', '');
                      //               end;
                    end
                  end;
                end;
              end;
            end;
          end
          else
          begin
            if ResStatus = 'Modify' then
            begin
              for Count2 := 0 to XMLRPHList.Count - 1 do
              begin
                DeployRPH(XMLRPHList.Strings[Count2], RoomStayID, RPH);
                if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                  DebugLogDetail.Add('InsertRV, RPH: ' + RPH);
                if RPH <> XMLEmptyData then
                begin
                  DeployRoomStay(XMLRoomStayList.Strings[RoomStayID], RoomTypeCode, BedTypeCode, RoomRateCode, Adult, Child, ArrivalDate, DepartureDate);
                  if BedTypeCode = '' then
                  begin
                    RoomList.Text := GetAvailableRoomByType(RoomTypeCode, '', FormatDateTimeX(ArrivalDate), FormatDateTimeX(DepartureDate), 0, 0, 0, 0, False, ProgramConfiguration.CCMSReservationAsAllotment);
                    if RoomList.Count > 0 then
                      BedTypeCode := GetBedTypeCode(RoomList.Strings[0]);
                  end;

                  if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                    DebugLogDetail.Add('DeployRoomStay: ' + ':' + RoomTypeCode + ':' + BedTypeCode + ':' + RoomRateCode + ':' + RoomRateAmountStr + ':' + IntToStr(Adult) + ':' + IntToStr(Child) + ':' + FormatDateTimeX2(ArrivalDate) + ':' + FormatDateTimeX2(DepartureDate));

                  for Count3 := 0 to XMLResGuestList.Count - 1 do
                  begin
                    DeployResGuest(XMLResGuestList.Strings[Count3], ArrivalTime, ResGuestRPH, FullName, Phone1, Email, Street, City, PostalCode, State, Country, Company);
                    if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
                      DebugLogDetail.Add('RPH = ResRPH ' + RPH + ':' + ResGuestRPH);
                    if RPH = ResGuestRPH then
                    begin
                      with DataModuleMain do
                      begin
                        try
                          if ProgramConfiguration.ChannelManagerVendor = ChannelManagerVendor.SiteMinder then
                            ChangeQueryString(MyQReservation,
                              'SELECT * FROM reservation' +
                              ' WHERE booking_code="' +BookingCode+ '"' +
                              ' AND booking_code<>"" ' +
                              'ORDER BY number;',
                              '', '', '', '', '', '', '', '', '', '')
                          else
                            ChangeQueryString(MyQReservation,
                              'SELECT * FROM reservation' +
                              ' WHERE booking_code="' +BookingCode+ '"' +
                              ' AND booking_code<>""' +
                              ' AND ota_id="' +OTAID+ '"' +
                              ' AND ota_id<>"" ' +
                              'ORDER BY number;',
                              '', '', '', '', '', '', '', '', '', '');
                        except
                        end;

                        if not MyQReservation.IsEmpty then
                        begin
                          ProgramConfiguration.CheckOutLimit := StrToTime(ReadConfigurationString(SystemCode.Hotel, ConfigurationCategory.Reservation, ConfigurationName.CheckOutLimit, False), ProgramVariable.FormatSettingX);
                          ReplaceTime(ArrivalDate, ArrivalTime);
                          ReplaceTime(DepartureDate, ProgramConfiguration.CheckOutLimit);

                          RoomRateAmount := 0;
                          if (RoomRateAmountStr <> '') and (RoomRateAmountStr <> '0') then
                          begin
                            try
                              RoomRateAmount := StrToFloat(RoomRateAmountStr, FormatStringTemp);
                            except
                            end;
                          end;

                          MyQReservation.First;
                          while not MyQReservation.Eof do
                          begin
                            ParameterCondition := MyQReservationnumber.AsLargeInt;
                            ContactPersonID := MyQReservationcontact_person_id.AsString;
                            GuestDetailID := MyQReservationguest_detail_id.AsString;
                            GuestProfileID := MyQReservationguest_profile_id.AsString;

                            UpdateContactPerson(ContactPersonID, '', FullName, Street, '', City, '', '', '', PostalCode, Phone1, '', '', Email, '', '', '', '', '', BoolToStr(True), '', '0000-00-00',
                                                '', '', '', '', '', '', '', '', '', '', '', '',
                                                '', '', '', '', '', '', '', '', '', '', '', '');
                            UpdateGuestDetail(GuestDetailID, RoomTypeCode, BedTypeCode, '', RoomRateCode, '', '', '', '', '', '', ArrivalDate, DepartureDate, Adult, Child, True, RoomRateAmount, RoomRateAmount, 0, 0);
                            UpdateGuestProfile(GuestProfileID, '', FullName, Street, '', City, '', '', '', PostalCode, Phone1, '', '', Email, '', '', '', '', '', BoolToStr(True), '',  '0000-00-00',
                                               '', '', '', '', '', '', '', '', '', '', '', '',
                                               '', '', '', '', '', '', '', '', '', '', '', '',
                                               ServerDate);

                            if (DateOf(ArrivalDate) >= DateOf(ProgramVariable.AuditDate)) and (DateOf(DepartureDate) > DateOf(ProgramVariable.AuditDate)) and
                               (GetAvailableRoomCountByType(RoomTypeCode, BedTypeCode, ArrivalDate, DepartureDate, ParameterCondition, 0, 0, 0, False, ProgramConfiguration.CCMSReservationAsAllotment) > 0) then
                            begin
                              UpdateReservation(ParameterCondition, '', '', '', '', GuestProfileID, '', '', '', '', FullName, '', '', ReservationStatus.New, '', '', '', '', '', '', OTAID, ResStatus, 1, NullDate, NullDate, True, False);

                              AssignRoom(ParameterCondition, False)
                            end
                            else
                              UpdateReservation(ParameterCondition, '', '', '', '', GuestProfileID, '', '', '', '', FullName, '', '', ReservationStatus.WaitList, '', '', '', '', '', '', OTAID, ResStatus, 1, NullDate, NullDate, True, False);
                              //                            InsertLogUser(LogUserAction.InsertReservation, IntToStr(ParameterCondition), '', '', '', LogUserAction.InsertReservationX);
                              //    //          ProcessSMSSchedule(SMSevent.OnInsertReservation, 'reservation.number = "' +IntToStr(ParameterCondition)+ '"', '', '', '', '', '', '', '');

                            UpdateReservationIsCMConfirmed(MyQReservationbooking_code.AsString, MyQReservationota_id.AsString, OTAID, 'Modify', False);

                            MyQReservation.Next;
                          end;
                        end;
                      end;
                    end;
                  end;
                  Break;
                end;
              end;
              Break;
            end
            else if ResStatus = 'Cancel' then
            begin
              with DataModuleMain do
              begin
                try
                  if ProgramConfiguration.ChannelManagerVendor = ChannelManagerVendor.SiteMinder then
                    ChangeQueryString(MyQReservation,
                      'SELECT * FROM reservation' +
                      ' WHERE booking_code="' +BookingCode+ '"' +
                      ' AND booking_code<>""' +
                      'ORDER BY number;',
                      '', '', '', '', '', '', '', '', '', '')
                  else
                    ChangeQueryString(MyQReservation,
                      'SELECT * FROM reservation' +
                      ' WHERE booking_code="' +BookingCode+ '"' +
                      ' AND booking_code<>""' +
                      ' AND ota_id="' +OTAID+ '"' +
                      ' AND ota_id<>"" ' +
                      'ORDER BY number;',
                      '', '', '', '', '', '', '', '', '', '');
                except
                end;

                if not MyQReservation.IsEmpty then
                begin
                  MyQReservation.First;
                  while not MyQReservation.Eof do
                  begin
                    ParameterCondition := MyQReservationnumber.AsLargeInt;

                    UpdateReservationStatus(ParameterCondition, ReservationStatus.Canceled, UserInfo.ID, 'Cancel by Channel Manager');
                    //                      InsertLogUser(LogUserAction.CancelReservation, IntToStr(ReservationNumber), '', '', Reason, LogUserAction.CancelReservationX);
                    UpdateReservationIsCMConfirmed(MyQReservationbooking_code.AsString, MyQReservationota_id.AsString, OTAID, 'Cancel', False);

                    MyQReservation.Next;
                  end;
                end;
              end;
            end;
          end;
        end;
      end;
    end;
  except
  end;

  if ProgramVariable.SQLDebugActive or ProgramVariable.DebugReservation then
  begin
    DebugLogDetail.SaveToFile('C:\Temp\XMLLogDetail' +IntToStr(ProgramVariable.DebugCount)+ '.txt');
    Inc(ProgramVariable.DebugCount);
    DebugLog.SaveToFile('C:\Temp\XMLLog' +IntToStr(ProgramVariable.DebugCount)+ '.txt');
    Inc(ProgramVariable.DebugCount);
  end;
end;